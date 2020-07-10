package goqtt

import (
	"fmt"
	"io"

	"github.com/rs/zerolog/log"

	"github.com/andschneider/goqtt/packets"
)

// Subscribe attempts to create a subscription for a client to it's configured topic.
// It sends a SUBSCRIBE packet and reads the SUBACK packet.
func (c *Client) Subscribe() error {
	// create packet
	var p packets.SubscribePacket
	p.CreatePacket(c.config.topic)

	c.stagePacket(&p)

	// read response and verify it's a SUBACK packet
	r, err := c.readResponse()
	if err != nil {
		return fmt.Errorf("could not read response for %s: %v", p.Name(), err)
	}
	if _, ok := r.(*packets.SubackPacket); !ok {
		typeErrorResponseLogger(p.Name(), r.Name(), r)
		return fmt.Errorf("did not receive a SUBACK packet, got %s instead", r.Name())
	}

	// start a KeepAlive process which will send Ping packets to prevent a disconnect
	// TODO I don't think this should be called in here - should be a background thing for a Client
	c.KeepAlive()
	return nil
}

// Unsubscribe sends an UNSUBSCRIBE packet for a given topic and reads the UNSUBACK packet.
func (c *Client) Unsubscribe() error {
	// create packet
	var p packets.UnsubscribePacket
	p.CreatePacket(c.config.topic)

	c.stagePacket(&p)

	// read response and verify it's a UNSUBACK packet
	r, err := c.readResponse()
	if err != nil {
		return fmt.Errorf("could not read response for %s: %v", p.Name(), err)
	}
	if _, ok := r.(*packets.UnsubackPacket); !ok {
		typeErrorResponseLogger(p.Name(), r.Name(), r)
		return fmt.Errorf("did not receive an UNSUBACK packet, got %s instead", r.Name())
	}
	return err
}

// TODO this is basically a wrapper around the ReadPacket function, and might get replaced/modified later.

// ReadLoop should be used after a successful subscription to a topic and reads any incoming messages.
// It returns a PublishPacket if one has been received for further processing.
func (c *Client) ReadLoop() (*packets.PublishPacket, error) {
	p, err := packets.ReadPacket(c.conn)
	if err != nil {
		if err == io.EOF {
			log.Warn().Msg("Looks like the server closed the connection...")
			return nil, err
		}
		log.Error().Err(err).Msg("subscribe loop error")
	}
	// process packets based on type
	switch packet := p.(type) {
	case *packets.PublishPacket:
		return packet, nil
	case *packets.PingRespPacket:
		// expected from the KeepAlive, all good
		log.Debug().Str("source", "goqtt").Str("packet", packet.String()).Msg("pingresp received")
		return nil, nil
	default:
		log.Warn().Str("source", "goqtt").Str("packet", packet.String()).Msg("packet type unexpected")
	}
	return nil, nil
}
