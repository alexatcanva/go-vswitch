package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/songgao/packets/ethernet"
)

var macTable map[string]string = make(map[string]string)

func main() {
	var (
		vswitchPort int
	)

	flag.IntVar(&vswitchPort, "port", 9999, "The port for the vswitch")

	l, err := net.ListenPacket("udp", fmt.Sprintf("0.0.0.0:%d", vswitchPort))
	if err != nil {
		log.Fatalf("unable to listen on udp port: %d, err: %s", vswitchPort, err.Error())
	}

	for {
		var frame ethernet.Frame
		frame.Resize(1518)
		_, remote, err := l.ReadFrom(frame)
		if err != nil {
			log.Printf("failed to read bytes into ethernet frame: %s\n", err.Error())
			continue
		}

		log.Printf("src %s -> dst %s\n", frame.Source().String(), frame.Destination().String())

		src := frame.Source().String()
		if v := macTable[src]; v != remote.String() {
			log.Printf("adding %s to table. %s -> %s\n", src, src, remote.String())
			macTable[src] = remote.String()
		}

		if dst, ok := macTable[src]; ok {
			u, err := net.ResolveUDPAddr("udp", dst)
			if err != nil {
				log.Printf("failed to resolve vswitch member: %s\n", err.Error())
				continue
			}
			forwardConn, err := net.DialUDP("udp", nil, u)
			if err != nil {
				log.Printf("failed to dial vswitch member: %s, err: %s\n", dst, err.Error())
				continue
			}
			_, err = forwardConn.Write(frame)
			if err != nil {
				log.Printf("failed to forward vswitch traffic to dst: %s, err: %s\n", dst, err.Error())
				continue
			}
		}
	}
}
