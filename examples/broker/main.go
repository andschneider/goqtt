/*
This is an example broker which allows clients to connect and publish messages. The broker keeps track of the clients
that are connected and which topics they are subscribed to. When a publish comes in from a client, the broker will send
along the message to each client subscribed to a topic.

to run: go run ./examples/broker/main.go

The broker listens on localhost and port 1884 by default. You can change the port with the -port flag.
*/

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"os"
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
	connecting  = make(chan client)
	leaving     = make(chan client)
	subscribe   = make(chan client)
	unsubscribe = make(chan client)
	messages    = make(chan packets.PublishPacket)
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
			log.Printf("client %q has connected.\n", cli.clientId)

		// subscribe clients to a topic
		case cli := <-subscribe:
			// create blank client map if not already present
			if _, ok := topics[cli.topic]; !ok {
				topics[cli.topic] = make(map[client]bool)
			}
			topics[cli.topic][cli] = true
			//fmt.Println(topics)

		// unsubscribe clients from a topic
		case cli := <-unsubscribe:
			for topic := range topics {
				delete(topics[topic], cli)
			}
			//fmt.Println(topics)

		case cli := <-leaving:
			log.Printf("got a request for a disconnect: %s\n", cli.clientId)
			delete(clients, cli)
			for topic := range topics {
				delete(topics[topic], cli)
			}
			// TODO remove empty topic maps
			fmt.Println(topics)
			// close client connection
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
			// stop timeout to prevent another disconnect request
			timer.Stop()
			//log.Println("Client disconnect channel is closed!")
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
			to := float64(cp.KeepAlive[1]) * 1.5 // timeout is 1.5 times the keep alive time
			cli = client{
				out:      c,
				timeout:  time.Duration(to) * time.Second,
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
			var ca packets.ConnackPacket
			ca.CreatePacket()
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
			var sa packets.SubackPacket
			sa.CreatePacket()
			err = sa.Write(c)
			if err != nil {
				log.Printf("could not send SUBACK packet: %v", err)
			}
		case packets.PingReqPacket:
			log.Printf("ping request received %v", p)
			// reset timeout
			timer.Reset(cli.timeout)

			// send pingresp packet
			var pp packets.PingRespPacket
			pp.CreatePacket()
			err = pp.Write(c)
			if err != nil {
				log.Printf("could not send PINGRESP packet: %v", err)
			}
		case packets.PublishPacket:
			var pRead, pWrite packets.PublishPacket
			log.Printf("publish received %v", p)
			// reset timeout
			timer.Reset(cli.timeout)

			// read publish packet
			pRead = p.(packets.PublishPacket)

			// send publish packet to be distributed to clients
			pWrite.CreatePacket(pRead.Topic, string(pRead.Message))
			messages <- pWrite

			// disconnect client after sending a message
			close(done) // close done channel to alert disconnect function
			leaving <- cli
		case packets.UnsubscribePacket:
			var u packets.UnsubackPacket
			log.Printf("unsubscribe request received %v", p)
			// reset timeout
			timer.Reset(cli.timeout)

			// send unsuback packet
			u.CreatePacket()
			err = u.Write(c)
			if err != nil {
				log.Printf("could not send UNSUBACK packet: %v", err)
			}

			// tell broker to remove client from subscription map
			unsubscribe <- cli
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
	server := flag.String("server", "127.0.0.1", "IP address to listen on. Default is localhost. Use 0.0.0.0 if running in Docker.")
	port := flag.String("port", "1884", "Port to allow connections on.")
	flag.Parse()

	if *port == "" {
		flag.Usage()
		os.Exit(1)
	}

	ln, err := net.Listen("tcp", *server+":"+*port)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening for clients on %s:%s", *server, *port)

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
