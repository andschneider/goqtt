package goqtt

import (
	"fmt"
	"io"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/andschneider/goqtt/packets"
)

// Subscribe attempts to create a subscription for a client to it's configured topic.
// It sends a SUBSCRIBE packet and reads the SUBACK packet.
func (c *Client) Subscribe() error {
	// create packet
	var p packets.SubscribePacket
	p.CreatePacket(c.config.topic)

	err := c.sendPacket(&p)
	if err != nil {
		return fmt.Errorf("could not send %s packet: %v", p.Name(), err)
	}

	// read response and verify it's a SUBACK packet
	r, err := c.readResponse()
	if err != nil {
		return fmt.Errorf("could not read response for %s: %v", p.Name(), err)
	}
	if _, ok := r.(*packets.SubackPacket); !ok {
		typeErrorResponseLogger(p.Name(), r.Name(), r)
		return fmt.Errorf("did not receive a SUBACK packet, got %s instead", r.Name())
	}
	return nil
}

// Unsubscribe sends an UNSUBSCRIBE packet for a given topic and reads the UNSUBACK packet.
func (c *Client) Unsubscribe() error {
	// create packet
	var p packets.UnsubscribePacket
	p.CreatePacket(c.config.topic)

	err := c.sendPacket(&p)
	if err != nil {
		return fmt.Errorf("could not send %s packet: %v", p.Name(), err)
	}

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

// SubscribeLoop keeps a connection alive after a successful subscription to a topic and reads any incoming messages.
// It sends pings based on the Keep Alive time to keep the connection from timing out.
func (c *Client) SubscribeLoop() {
	ticker := time.NewTicker(time.Duration(c.config.keepAlive) * time.Second)
	// TODO add disconnect functionality
	disconnect := make(chan bool)
	go func() {
		for {
			select {
			case <-disconnect:
				return
			case <-ticker.C:
				err := c.SendPing()
				if err != nil {
					log.Fatal().Err(err)
				}
				//fmt.Println("would send a ping")
			}
		}
	}()

	for {
		p, err := packets.ReadPacket(c.conn)
		// process packets based on type
		switch packet := p.(type) {
		case *packets.PublishPacket:
			log.Info().
				Str("TOPIC", packet.Topic).
				Str("DATA", string(packet.Message)).
				Msg("publish packet received")
		}
		if err != nil {
			if err == io.EOF {
				log.Warn().Msg("Looks like the server closed the connection...")
				break
			}
			log.Fatal().Err(err).Msg("subscribe loop error")
			return
		}
	}
}
