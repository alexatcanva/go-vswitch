package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"

	"github.com/alexbathome/go-vswitch/pkg/vport"
)

func main() {
	var (
		serverAddr string
		tapName    string
	)

	flag.StringVar(&serverAddr, "address", "", "The server address <ip>:<port>")
	flag.StringVar(&tapName, "tap", "gotap", "The name of the TAP device")
	flag.Parse()

	if serverAddr == "" {
		flag.Usage()
		os.Exit(1)
	}

	ifaceTap, err := vport.New(tapName)
	if err != nil {
		log.Fatalf("unable to create TAP iface: %s", err.Error())
	}

	udpServerAddress, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		log.Fatalf("unable to resolve server address: %s", err.Error())
	}

	conn, err := net.DialUDP("udp4", nil, udpServerAddress)
	if err != nil {
		log.Fatalf("unable to dial udp server: %s", err.Error())
	}

	for {
		go io.Copy(ifaceTap, conn)
		io.Copy(conn, ifaceTap)
	}
}
