// Package goqtt is a MQTT 3.1.1 client library.
// It does not implement the full client specification, see README for more information.
package goqtt

import (
	"fmt"
	"net"

	"github.com/andschneider/goqtt/packets"
	"github.com/rs/zerolog/log"
)

// ClientConfig holds the information needed to create a new Client
type ClientConfig struct {
	// ClientId is the identifier you use to tell the MQTT broker who you are.
	// A broker will require unique ClientId's for all connected clients.
	ClientId string
	// KeepAlive is the time in seconds that the broker should wait before disconnecting
	// you if no packets are sent.
	KeepAlive int
	// Server is the address (IP or domain) of the broker.
	Server string
	// Port is the port number of the broker. Usually 1883 (insecure) or 8883 (TLS).
	Port string
	// Topic is the full qualified topic to subscribe or publish to.
	// Only single topics and no wildcards are accepted at this time.
	Topic string
}

// Client is the main interaction point for sending and receiving Packets. Using the configuration
// set in the ClientConfig struct, an instantiated Client needs to call the Connect() method
// before sending/receiving any packets.
type Client struct {
	config *ClientConfig
	conn   net.Conn

	send chan packets.Packet
}

// NewClient creates a new Client based on the configuration values in the ClientConfig struct.
func NewClient(config *ClientConfig) (*Client, error) {
	if config.Server == "" {
		return nil, fmt.Errorf("ClientConfig: Server must be set")
	}
	if config.Port == "" {
		return nil, fmt.Errorf("ClientConfig: Port must be set")
	}
	return &Client{config: config}, nil
}

// Connect attempts to create a TCP connection to the broker specified in the client's ClientConfig.
// It sends a CONNECT packet and reads the CONNACK packet.
func (c *Client) Connect() error {
	// connect over TCP
	b := fmt.Sprintf("%s:%s", c.config.Server, c.config.Port)
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
	cp.KeepAlive = []byte{0, byte(c.config.KeepAlive)}
	cp.ClientIdentifier = c.config.ClientId

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
