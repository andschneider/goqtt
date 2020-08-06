package goqtt

import (
	"fmt"

	"github.com/andschneider/goqtt/packets"
)

// SendPublish sends a given message to the topic specified by the Client.
func (c *Client) SendPublish(message string) error {
	// create packet
	var p packets.PublishPacket
	p.CreatePublishPacket(c.config.topic, message)

	// TODO review this
	// Write directly to connection instead of sending packet to client's connection channel to avoid
	// race conditions if the SendPublish command happens towards the end of a script.
	err := p.Write(c.conn)
	if err != nil {
		return fmt.Errorf("could not write Publish packet: %v", err)
	}

	// has no response with QOS 0
	return nil
}
