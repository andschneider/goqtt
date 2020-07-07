package goqtt

import (
	"github.com/andschneider/goqtt/packets"
)

// SendPublish sends a given message to the topic specified by the Client.
func (c *Client) SendPublish(message string) error {
	// create packet
	var p packets.PublishPacket
	p.CreatePacket(c.config.topic, message)

	//err := c.sendPacket(&p)
	c.send <- &p
	//if err != nil {
	//	return fmt.Errorf("could not send %s packet: %v", p.Name(), err)
	//}

	// has no response with QOS 0
	return nil
}
