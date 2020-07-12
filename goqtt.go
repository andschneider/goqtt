package goqtt

import (
	"fmt"
	"math/rand"
	"time"

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
// data from a client's TCP connection. It basically wraps the ReadPacket with some
// helpful logging messages.
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

// stagePacket is a helper function which just sends the packet to the client channel.
// It will probable get removed or heavily modified.
func (c *Client) stagePacket(p packets.Packet) int {
	// create traceId for the packet, which is probably unnecessary.
	rand.Seed(time.Now().UnixNano())
	traceId := rand.Int()
	log.Debug().
		Str("source", "goqtt").
		Str("packetType", p.Name()).
		Str("packet", p.String()).
		Int("traceId", traceId).
		Msg("stage packet")
	c.send <- p
	return traceId
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
