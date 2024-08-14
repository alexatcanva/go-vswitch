package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/songgao/packets/ethernet"
	"golang.org/x/exp/maps"
)

var macTable map[string]*net.UDPAddr = make(map[string]*net.UDPAddr)

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
		n, remote, err := l.ReadFrom(frame)
		if err != nil {
			log.Printf("failed to read bytes into ethernet frame: %s\n", err.Error())
			continue
		}
		frame.Resize(n)

		log.Printf("src %s -> dst %s\n", frame.Source().String(), frame.Destination().String())

		src := frame.Source().String()
		if v := macTable[src]; v.String() != remote.(*net.UDPAddr).String() {
			log.Printf("adding %s to table. %s -> %s\n", src, src, remote)
			macTable[src] = remote.(*net.UDPAddr)
		}

		if dst, ok := macTable[frame.Destination().String()]; ok {
			_, err = l.WriteTo(frame, dst)
			log.Printf("forwarding vswitch traffic to dst: %s\n", dst)
			if err != nil {
				log.Printf("failed to forward vswitch traffic to dst: %s, err: %s\n", dst, err.Error())
				continue
			}
			continue
		}

		if frame.Destination().String() == "ff:ff:ff:ff:ff:ff" {
			for _, dst := range maps.Keys(macTable) {
				if dst == frame.Source().String() {
					log.Printf("ignoring broadcast back to host.\n")
					continue
				}
				log.Println("broadcasting to", dst)
				_, err = l.WriteTo(frame, macTable[dst])
				if err != nil {
					log.Printf("failed to forward vswitch traffic to dst: %s, err: %s\n", dst, err.Error())
					continue
				}
			}
		}
	}
}
