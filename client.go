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

// NewClientConfig creates a ClientConfig struct and since none of the methods are exported, they
// can't be changed after (this is on purpose)
func NewClientConfig(clientId string, keepAlive int, broker string, topic string) *ClientConfig {
	return &ClientConfig{clientId: clientId, keepAlive: keepAlive, broker: broker, topic: topic}
}

// Client is the main interaction point for sending and receiving Packets. Using the configuration
// set in the ClientConfig struct, an instantiated Client needs to call the Connect() method
// before sending/receiving any packets
type Client struct {
	config *ClientConfig
	conn   net.Conn
}

// NewClient creates a new Client based on the configuration values in the ClientConfig struct
func NewClient(config *ClientConfig) *Client {
	return &Client{config: config}
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

	// create Connect packet
	var cp packets.ConnectPacket
	cp.CreatePacket()
	cp.KeepAlive = []byte{0, byte(c.config.keepAlive)}
	cp.ClientIdentifier = c.config.clientId

	// send packet
	err = c.sendPacket(&cp)
	if err != nil {
		return fmt.Errorf("could not send %s packet: %v", cp.Name(), err)
	}

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

// Disconnect sends a DISCONNECT packet.
func (c *Client) Disconnect() error {
	// create packet
	var p packets.DisconnectPacket
	p.CreatePacket()

	err := c.sendPacket(&p)
	if err != nil {
		return fmt.Errorf("could not send %s packet: %v", p.Name(), err)
	}

	return nil
}

// Close closes the Client's TCP connection
func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		fmt.Printf("could not close Client connection: %v", err)
	}
}
