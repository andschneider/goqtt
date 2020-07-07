package goqtt

import (
	"fmt"

	"github.com/andschneider/goqtt/packets"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// disable logging by default. over ride by calling function again with a different log level, like below:
	//zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

// readResponse is a helper method to the ReadPacket function which attempts to read
// data from a TCP connection.
func (c *Client) readResponse() (packets.Packet, error) {
	p, err := packets.ReadPacket(c.conn)
	if err != nil {
		log.Error().
			Err(err).
			Str("source", "goqtt").
			Msg("could not read a Packet from the TCP connection")
		return nil, fmt.Errorf("could not read a Packet from the TCP connection: %v", err)
	}
	log.Debug().
		Str("source", "goqtt").
		Str("packetType", p.Name()).
		Str("packet", p.String()).
		Msg("receiving packet")
	return p, nil
}

// sendPacket is a helper function to write a Packet to a TCP connection
func (c *Client) sendPacket(p packets.Packet) error {
	//log.Debug().
	//	Str("source", "goqtt").
	//	Str("packetType", p.Name()).
	//	Str("packet", p.String()).
	//	Msg("sending packet")
	c.send <- p
	//err := p.Write(c.conn)
	//if err != nil {
	//	log.Error().
	//		Err(err).
	//		Str("source", "goqtt").
	//		Str("packetType", p.Name()).
	//		Str("packet", p.String()).
	//		Msg("could not send packet")
	//	return fmt.Errorf("could not send Packet to TCP connection: %v", err)
	//}
	return nil
}

// typeErrorResponseLogger is a helper that logs relevant information when the wrong type
// of packet is received from the TCP connection.
func typeErrorResponseLogger(sendType, receiveType string, packet packets.Packet) {
	log.Error().
		Str("source", "goqtt").
		Str("sentType", sendType).
		Str("receivedType", receiveType).
		Str("packet", packet.String()).
		Msg("received wrong type of packet")
}
