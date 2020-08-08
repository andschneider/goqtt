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
