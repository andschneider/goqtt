package goqtt

import (
	"fmt"
	"net"

	"github.com/andschneider/goqtt/packets"
	"github.com/rs/zerolog/log"
)

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

	// create Connect packet
	var cp packets.ConnectPacket
	cp.CreatePacket()
	cp.KeepAlive = []byte{0, byte(c.Config.keepAlive)}
	cp.ClientIdentifier = c.Config.clientId

	// send packet
	err = c.sendPacket(&cp)
	if err != nil {
		return fmt.Errorf("could not write Connect packet: %v", err)
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
func (c *Client) Disconnect() {
	// create packet
	var p packets.DisconnectPacket
	p.CreatePacket()

	err := c.sendPacket(&p)
	if err != nil {
		log.Error().Err(err).Str("source", "goqtt").Msg("could not write Disconnect packet")
	}
	err = c.conn.Close()
	if err != nil {
		log.Error().Err(err).Str("source", "goqtt").Msg("could not close Client connection")
	}
}
