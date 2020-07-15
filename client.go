package goqtt

import (
	"fmt"
	"net"

	"github.com/andschneider/goqtt/packets"
	"github.com/rs/zerolog/log"
)

// ClientConfig holds the information needed to create a new Client
type ClientConfig struct {
	clientId  string
	keepAlive int
	broker    string // address and port
	topic     string
}

// NewClientConfig creates a ClientConfig struct which can't be edited after it's created.
func NewClientConfig(clientId string, keepAlive int, broker string, topic string) *ClientConfig {
	return &ClientConfig{clientId: clientId, keepAlive: keepAlive, broker: broker, topic: topic}
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
func NewClient(config *ClientConfig) *Client {
	return &Client{config: config}
}

// TODO might combine these Get helpers

// GetClientId returns the current ClientId for the client.
func (c *Client) GetClientId() string {
	return c.config.clientId
}

// GetKeepAlive returns the current KeepAlive for the client.
func (c *Client) GetKeepAlive() int {
	return c.config.keepAlive
}

// GetBroker returns the current Broker for the client.
func (c *Client) GetBroker() string {
	return c.config.broker
}

// GetTopic returns the current Topic for the client.
func (c *Client) GetTopic() string {
	return c.config.topic
}

// Connect attempts to create a TCP connection to the broker specified in the client's ClientConfig.
// It sends a CONNECT packet and reads the CONNACK packet.
func (c *Client) Connect() error {
	// connect over TCP
	conn, err := net.Dial("tcp", c.config.broker)
	if err != nil {
		log.Error().Err(err).
			Str("source", "goqtt").
			Str("broker", c.config.broker).
			Msg("could not connect to server")
		return fmt.Errorf("could not connect to server: %v", err)
	}
	c.conn = conn

	c.send = make(chan packets.Packet)
	go c.sendPackets()

	// create Connect packet
	var cp packets.ConnectPacket
	cp.CreatePacket()
	cp.KeepAlive = []byte{0, byte(c.config.keepAlive)}
	cp.ClientIdentifier = c.config.clientId

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

// sendPacket handles all the packets sent through the client's channel and writes them to the TCP connection
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
