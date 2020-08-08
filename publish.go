package goqtt

import (
	"fmt"

	"github.com/andschneider/goqtt/packets"
)

// Publish sends a given message to the topic specified by the Client.
func (c *Client) Publish(message string) error {
	// create packet
	var p packets.PublishPacket
	p.CreatePublishPacket(c.Config.topic, message)

	err := c.sendPacket(&p)
	if err != nil {
		return fmt.Errorf("could not write Publish packet: %v", err)
	}

	// has no response with QOS 0
	return nil
}
