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
	"io"
	"math"
	"net"
	"os"
	"strings"
	"time"

	"github.com/andschneider/goqtt/packets"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
				log.Debug().
					Str("clientId", cli.clientId).
					Str("packet", msg.String()).
					Msg("attempting to send Publish packet")
				err := msg.Write(cli.out)
				if err != nil {
					log.Error().Err(err).Str("clientId", cli.clientId).Msg("could not write publish packet")
				}
				log.Info().
					Str("clientId", cli.clientId).
					Str("TOPIC", msg.Topic).
					Str("DATA", string(msg.Message)).
					Msg("sent publish packet")
			}

		// register new clients
		case cli := <-connecting:
			log.Debug().
				Str("clientId", cli.clientId).
				Str("request", "connect").
				Msg("got a request for a connection")
			clients[cli] = true
			log.Info().
				Str("clientId", cli.clientId).
				Str("status", "connected").
				Msg("client has connected successfully.")

		// subscribe clients to a topic
		case cli := <-subscribe:
			log.Debug().
				Str("clientId", cli.clientId).
				Str("request", "subscribe").
				Msg("got a request for a subscribe")
			// create blank client map if not already present
			if _, ok := topics[cli.topic]; !ok {
				topics[cli.topic] = make(map[client]bool)
			}
			topics[cli.topic][cli] = true
			log.Info().
				Str("TOPIC", cli.topic).
				Str("clientId", cli.clientId).
				Msg("subscribe successful.")

		// unsubscribe clients from a topic
		case cli := <-unsubscribe:
			log.Debug().
				Str("clientId", cli.clientId).
				Str("request", "unsubscribe").
				Msg("got a request for an unsubscribe")
			for topic := range topics {
				delete(topics[topic], cli)
				log.Info().
					Str("TOPIC", cli.topic).
					Str("clientId", cli.clientId).
					Msg("unsubscribe successful.")
			}

		case cli := <-leaving:
			log.Debug().
				Str("clientId", cli.clientId).
				Str("request", "disconnect").
				Msg("got a request for a disconnect")
			delete(clients, cli)
			for topic := range topics {
				delete(topics[topic], cli)
			}
			// TODO remove empty topic maps
			//fmt.Println(topics)
			// close client connection
			err := cli.out.Close()
			if err != nil {
				log.Error().
					Err(err).
					Str("clientId", cli.clientId).
					Msgf("error closing network connection for %v", cli)
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
			log.Debug().Msg("Client disconnect channel is closed!")
			break
		}
		p, err := packets.ReadPacket(c)
		if err != nil {
			if err == io.EOF {
				// TODO do i care about EOFs?
				log.Warn().Msg("received an EOF")
				break
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
			log.Error().Err(err).Msg("unknown ReadPacket error")
		}

		switch t := p.(type) {
		// try to read connection packet first
		// what if it's not a connection packet for a new client?
		case *packets.ConnectPacket:
			cli.packetTrace(p)
			// read in connection information and register new client with broker
			cp := p.(*packets.ConnectPacket)
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
				log.Warn().
					Str("clientId", cli.clientId).
					Str("op", "timeout").
					Msg("client timed out")
				leaving <- cli
			}()

			// send a connack
			var ca packets.ConnackPacket
			ca.CreatePacket()
			err = ca.Write(c)
			if err != nil {
				cli.packetError(&ca, err)
			}
		case *packets.SubscribePacket:
			cli.packetTrace(p)
			// read subscribe packet
			sp := p.(*packets.SubscribePacket)
			for _, t := range sp.Topics {
				cli.topic = t
				subscribe <- cli // send topic info
			}

			// send suback packet
			var sa packets.SubackPacket
			sa.CreatePacket()
			err = sa.Write(c)
			if err != nil {
				cli.packetError(&sa, err)
			}
		case *packets.PingReqPacket:
			cli.packetTrace(p)
			// reset timeout
			timer.Reset(cli.timeout)

			// send pingresp packet
			var pp packets.PingRespPacket
			pp.CreatePacket()
			err = pp.Write(c)
			if err != nil {
				cli.packetError(&pp, err)
			}
		case *packets.PublishPacket:
			cli.packetTrace(p)
			// reset timeout
			timer.Reset(cli.timeout)

			// read publish packet
			pRead := p.(*packets.PublishPacket)

			// send publish packet to be distributed to clients
			var pWrite packets.PublishPacket
			pWrite.CreatePacket(pRead.Topic, string(pRead.Message))
			messages <- pWrite

			// disconnect client after sending a message
			close(done) // close done channel to alert disconnect function
			leaving <- cli
		case *packets.UnsubscribePacket:
			cli.packetTrace(p)
			//log.Debug().Str("packetType", "Unsubscribe").Str("packet", p.String()).Msg("packet received")
			// reset timeout
			timer.Reset(cli.timeout)

			// send unsuback packet
			var u packets.UnsubackPacket
			u.CreatePacket()
			err = u.Write(c)
			if err != nil {
				cli.packetError(&u, err)
			}

			// tell broker to remove client from subscription map
			unsubscribe <- cli
		case *packets.DisconnectPacket:
			cli.packetTrace(p)
			log.Info().Str("clientId", cli.clientId).Msg("disconnect received")
			close(done) // close done channel to alert disconnect function
			leaving <- cli
		default:
			if t == nil {
				return
			}
			log.Warn().Str("op", "PacketRead").Msgf("unexpected packet type %t", t)
			return
		}
	}
}

func (c *client) packetTrace(p packets.Packet) {
	log.Trace().
		Str("clientId", c.clientId).
		Str("packetType", p.Name()).
		Str("packet", p.String()).
		Msg("packet received")
}

func (c *client) packetError(p packets.Packet, err error) {
	log.Error().
		Err(err).
		Str("clientId", c.clientId).
		Str("packetType", p.Name()).
		Str("packet", p.String()).
		Msg("could not send packet")
}

func main() {
	server := flag.String("server", "127.0.0.1", "IP address to listen on. Default is localhost. Use 0.0.0.0 if running in Docker.")
	port := flag.String("port", "1884", "Port to allow connections on.")
	verbose := flag.Bool("v", false, "Verbose output. Default is false.")
	extraVerbose := flag.Bool("vv", false, "Extra verbose output. Will override the verbose flag. Default is false.")
	flag.Parse()

	// Set logger to pretty print instead of structured json
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Set log level to debug if verbose is passed in
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if *extraVerbose {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	if *port == "" {
		flag.Usage()
		os.Exit(1)
	}

	ln, err := net.Listen("tcp", *server+":"+*port)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msgf("Listening for clients on %s:%s", *server, *port)

	go broker()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Error().Err(err).Msg("error on client connection.")
			continue
		}
		log.Info().Msg("client connecting...")
		go handleConnection(conn)
	}
}
