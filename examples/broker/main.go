package main

import (
	"fmt"
	"log"
	"net"

	"github.com/andschneider/goqtt/packets"
)

type client struct {
	out     chan<- []byte
	timeout int64
	//topics   []string // TODO can't use as map key
	topic    string
	clientId string
}

type publish struct {
	message []byte
	topic   string
}

var (
	connecting = make(chan client)
	leaving    = make(chan client)
	messages   = make(chan publish)
	subscribe  = make(chan client)
)

func broker() {
	clients := make(map[client]bool)    // all connected clients
	topics := make(map[string][]client) // TODO is this easier for publishing based on topics?
	for {
		select {
		case msg := <-messages:
			// Send messages to every client subscribed to a topic
			for _, cli := range topics[msg.topic] {
				cli.out <- msg.message
			}

		case cli := <-connecting:
			clients[cli] = true

		case cli := <-subscribe:
			topics[cli.topic] = append(topics[cli.topic], cli)
			fmt.Println(topics)

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.out)
		}
	}
}

func handleConnection(c net.Conn) {
	var cli client
	ch := make(chan []byte) // outgoing messages
	go packetWriter(c, ch)

	for {
		p, err := packets.Reader(c)
		if err != nil {
			log.Print(err)
		}
		switch t := p.(type) {
		// try to read connection packet first
		// what if it's not a connection packet for a new client?
		case packets.ConnectPacket:
			log.Printf("connect packet recieved %v", p)
			// read in connection information and register new client with broker
			cp := p.(packets.ConnectPacket)
			cli = client{
				out:      ch,
				timeout:  int64(cp.KeepAlive[1]),
				clientId: cp.ClientIdentifier,
			}
			connecting <- cli // register client with broker

			// send a connack
			ca := packets.CreateConnackPacket()
			err = ca.Write(c)
			if err != nil {
				log.Printf("could not send CONNACK packet: %v", err)
			}

			// TODO not sure if want to send it over the channel
			//buf := new(bytes.Buffer)
			//err := ca.Write(buf)
			//if err != nil {
			//	log.Printf("could not write connack packet: %v", err)
			//}
			//ch <- buf.Bytes()
		case packets.SubscribePacket:
			log.Printf("subscribe packet recieved %v", p)
			// read subscribe packet
			sp := p.(packets.SubscribePacket)
			for _, t := range sp.Topics {
				cli.topic = t
				subscribe <- cli // send topic info
			}

			// send suback packet
			sap := packets.CreateSubackPacket()
			err = sap.Write(c)
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
		case packets.PublishPacket:
			log.Printf("publish received %v", p)

			// TODO send publish packet
			//buf := new(bytes.Buffer)
			//pp := packets.CreatePublishPacket()
			//err := pp.Write(buf)
			//if err != nil {
			//	log.Printf("could not write connack packet: %v", err)
			//}
			//ch <- buf.Bytes()
		default:
			fmt.Printf("unexpected type %t\n", t)
		}
		if err != nil {
			log.Println("client left..")
			c.Close()
			// escape recursion
			return
		}
	}
}

func packetWriter(conn net.Conn, ch <-chan []byte) {
	for msg := range ch {
		_, err := conn.Write(msg)
		if err != nil {
			log.Printf("could not send packet: %v", err)
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:1884")
	if err != nil {
		log.Fatal(err)
	}

	go broker()
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
