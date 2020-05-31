package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"strings"
	"time"

	"github.com/andschneider/goqtt/packets"
)

type client struct {
	out     net.Conn
	timeout time.Duration
	//topics   *[]string
	topic    string
	clientId string
}

var (
	connecting = make(chan client)
	leaving    = make(chan client)
	subscribe  = make(chan client)
	messages   = make(chan packets.PublishPacket)
)

func broker() {
	clients := make(map[client]bool) // all connected clients
	topics := make(map[string]map[client]bool)
	for {
		select {
		case msg := <-messages:
			// Send messages to every client subscribed to a topic
			for cli := range topics[msg.Topic] {
				//fmt.Printf("trying to send  %s to %s\n", msg.Message, msg.Topic)
				err := msg.Write(cli.out)
				if err != nil {
					log.Printf("could not write publish packet: %v", err)
				}
			}

		// register new clients
		case cli := <-connecting:
			clients[cli] = true

		// subscribe clients to a topic
		case cli := <-subscribe:
			// create blank client map if not already present
			if _, ok := topics[cli.topic]; !ok {
				topics[cli.topic] = make(map[client]bool)
			}
			topics[cli.topic][cli] = true
			//fmt.Println(topics)

		case cli := <-leaving:
			log.Printf("got a request for a disconnect: %s\n", cli.clientId)
			delete(clients, cli)
			for topic := range topics {
				delete(topics[topic], cli)
			}
			// TODO remove empty topic maps
			fmt.Println(topics)
			err := cli.out.Close()
			if err != nil {
				log.Printf("error closing network connection for %v: %v\n", cli, err)
			}
		}
	}
}

// disconnected checks whether the channel has been closed
func disconnected(ch <-chan bool) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

func handleConnection(c net.Conn) {
	var cli client
	// initialize timer
	timer := time.NewTimer(math.MaxInt64)
	timer.Stop()
	// disconnect channel
	done := make(chan bool)

	for {
		if disconnected(done) {
			// stop timeout to prevent another disconnect request and break out of loop
			//log.Println("Client disconnect channel is closed!")
			timer.Stop()
			break
		}
		p, err := packets.Reader(c)
		if err != nil {
			if err == io.EOF {
				// TODO do i care about EOFs?
				break
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
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
				out:      c,
				timeout:  time.Duration(cp.KeepAlive[1]) * time.Second,
				clientId: cp.ClientIdentifier,
			}
			connecting <- cli // register client with broker

			// timeout
			//timer.Reset(10 * time.Second)
			timer.Reset(cli.timeout)
			go func() {
				<-timer.C
				log.Printf("client %s timed out\n", cli.clientId)
				leaving <- cli
			}()

			// send a connack
			ca := packets.CreateConnackPacket()
			err = ca.Write(c)
			if err != nil {
				log.Printf("could not send CONNACK packet: %v", err)
			}
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
			// reset timeout
			timer.Reset(cli.timeout)

			// send pingresp packet
			pp := packets.CreatePingRespPacket()
			err = pp.Write(c)
			if err != nil {
				log.Printf("could not send PINGRESP packet: %v", err)
			}
		case packets.PublishPacket:
			log.Printf("publish received %v", p)
			// read publish packet
			pp := p.(packets.PublishPacket)

			// send publish packet to be distributed to clients
			ppp := packets.CreatePublishPacket(pp.Topic, string(pp.Message))
			messages <- ppp
		case packets.DisconnectPacket:
			log.Printf("disconnect received %v", p)
			close(done) // close done channel to alert disconnect function
			leaving <- cli
		default:
			if t == nil {
				return
			}
			fmt.Printf("unexpected type %t\n", t)
			return
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
		log.Println("client connecting...")
		go handleConnection(conn)
	}
}
