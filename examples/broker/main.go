package main

import (
	"fmt"
	"log"
	"net"

	"github.com/andschneider/goqtt/packets"
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:1884")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Println("asdf")
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	// try to read connection packet first
	p, err := packets.Reader(c)
	if err != nil {
		log.Print(err)
	}
	switch t := p.(type) {
	case packets.ConnectPacket:
		log.Printf("connect packet recieved %v", p)
		// do something like send a connack and save connection information
		cp := packets.CreateConnackPacket()
		err = cp.Write(c)
		if err != nil {
			log.Printf("could not send CONNACK packet: %v", err)
		}
	case packets.SubscribePacket:
		log.Printf("subscribe packet recieved %v", p)
		sp := packets.CreateSubackPacket()
		err = sp.Write(c)
		if err != nil {
			log.Printf("could not send SUBACK packet: %v", err)
		}
	case packets.PingReqPacket:
		log.Printf("ping request received %v", p)
		pp := packets.CreatePingRespPacket()
		err = pp.Write(c)
		if err != nil {
			log.Printf("could not send PINGRESP packet: %v", err)
		}
	default:
		fmt.Printf("unexpected type %t\n", t)
	}
	if err != nil {
		log.Println("client left..")
		c.Close()
		// escape recursion
		return
	}

	// TODO I think there's a better way to do this. maybe channels?
	fmt.Println("recursing")
	handleConnection(c)
}
