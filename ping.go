package goqtt

import (
	"fmt"

	"github.com/andschneider/goqtt/packets"
	"github.com/rs/zerolog/log"
)

// SendPing is a helper function to create a PINGREQ packet and send it right away.
// It also reads the PINGRESP packet.
func (c *Client) SendPing() error {
	// create packet
	var p packets.PingReqPacket
	p.CreatePacket()

	err := c.sendPacket(&p)
	if err != nil {
		return fmt.Errorf("could not send %s packet: %v", p.Name(), err)
	}

	// response
	// why did i make this a goroutine?
	go func() {
		_, err := c.readResponse()
		if err != nil {
			log.Fatal().Err(err)
		}
	}()
	return nil
}
