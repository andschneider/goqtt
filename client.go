// Package goqtt is a MQTT 3.1.1 client library.
// It does not implement the full client specification, see README for more information.
package goqtt

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/andschneider/goqtt/packets"
	"github.com/rs/zerolog/log"
)

// clientConfig holds the information needed to create a new Client.
// As it's unexported these can't be set directly, instead use a functional
// option, like Port.
type clientConfig struct {
	clientId  string
	keepAlive int
	// server is the address (IP or domain) of the broker.
	server string
	port   string
	topic  string
}

// Client is the main interaction point for sending and receiving Packets. An
// instantiated Client needs to call the Connect() method before sending/receiving any packets.
type Client struct {
	Config *clientConfig
	conn   net.Conn

	send chan packets.Packet
}

// NewClient creates a Client struct which can be used to interact with a MQTT broker.
// It sets default values for most of the configuration, so only a server address is
// required to instantiate it. Other configuration options are available and can be
// passed in as needed.
func NewClient(addr string, opts ...option) *Client {
	var c = &Client{}
	// create a default ClientId based on time to reduce collisions in the broker
	cid := fmt.Sprintf("%s-%d-%s", "goqtt", os.Getpid(), strconv.Itoa(time.Now().Second()))
	// default configuration
	d := &clientConfig{
		clientId:  cid,
		keepAlive: 60,
		server:    addr,
		port:      "1883",
		topic:     "goqtt",
	}
	c.Config = d
	// apply any other options that have been passed in
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// option is a function used to modify/set a configuration field in a Client.
type option func(*Client)

// ClientId is the identifier you use to tell the MQTT broker who you are.
// A broker will require unique ClientId's for all connected clients.
func ClientId(cid string) option {
	return func(c *Client) {
		c.Config.clientId = cid
	}
}

// KeepAlive is the time in seconds that the broker should wait before
// disconnecting you if no packets are sent.
func KeepAlive(seconds int) option {
	return func(c *Client) {
		c.Config.keepAlive = seconds
	}
}

// Port is the port number of the server. Usually 1883 (insecure) or 8883 (TLS).
func Port(port string) option {
	return func(c *Client) {
		c.Config.port = port
	}
}

// Topic is the fully qualified topic to subscribe or publish to.
// Only single topics and no wildcards are accepted at this time.
func Topic(topic string) option {
	return func(c *Client) {
		c.Config.topic = topic
	}
}

// Connect attempts to create a TCP connection to the broker specified in the client's ClientConfig.
// It sends a CONNECT packet and reads the CONNACK packet.
func (c *Client) Connect() error {
	// connect over TCP
	b := fmt.Sprintf("%s:%s", c.Config.server, c.Config.port)
	conn, err := net.Dial("tcp", b)
	if err != nil {
		log.Error().Err(err).
			Str("source", "goqtt").
			Str("broker", b).
			Msg("could not connect to server")
		return fmt.Errorf("could not connect to server: %v", err)
	}
	c.conn = conn

	c.send = make(chan packets.Packet)
	go c.sendPackets()

	// create Connect packet
	var cp packets.ConnectPacket
	cp.CreatePacket()
	cp.KeepAlive = []byte{0, byte(c.Config.keepAlive)}
	cp.ClientIdentifier = c.Config.clientId

	// send packet
	c.stagePacket(&cp)

	// read response and verify it's a CONNACK packet
	r, err := c.readResponse()
	if err != nil {
		return fmt.Errorf("could not read response for %s: %v", cp.Name(), err)
	}
	if _, ok := r.(*packets.ConnackPacket); !ok {
		typeErrorResponseLogger(cp.Name(), r.Name(), r)
		return fmt.Errorf("did not receive a CONNACK packet, got %s instead", r.Name())
	}
	return nil
}

// sendPackets handles all the packets sent through the client's channel and writes them to the TCP connection
func (c *Client) sendPackets() {
	for {
		p := <-c.send
		log.Debug().
			Str("source", "goqtt").
			Str("packetType", p.Name()).
			Str("packet", p.String()).
			Msg("send packet")
		err := p.Write(c.conn)
		if err != nil {
			log.Error().
				Err(err).
				Str("source", "goqtt").
				Str("packetType", p.Name()).
				Str("packet", p.String()).
				Msg("could not send packet")
		}
	}
}

// Disconnect sends a DISCONNECT packet.
func (c *Client) Disconnect() {
	// create packet
	var p packets.DisconnectPacket
	p.CreatePacket()

	err := p.Write(c.conn)
	if err != nil {
		log.Error().Err(err).Str("source", "goqtt").Msg("could not write Disconnect packet")
	}
	err = c.conn.Close()
	if err != nil {
		log.Error().Err(err).Str("source", "goqtt").Msg("could not close Client connection")
	}
}
