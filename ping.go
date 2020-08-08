package goqtt

import (
	"fmt"
	"time"

	"github.com/andschneider/goqtt/packets"
	"github.com/rs/zerolog/log"
)

// SendPing creates a Ping packet and writes it to the client's TCP connection right away.
func (c *Client) SendPing() error {
	// create packet
	var p packets.PingReqPacket
	p.CreatePacket()

	err := p.Write(c.conn)
	if err != nil {
		return fmt.Errorf("could not write Ping packet: %v", err)
	}
	return nil
}

// keepAlive is an unsophisticated way to prevent the server from closing the client's connection.
// It blindly sends a Ping packet according to the client's configured KeepAlive setting. This
// is possibly wasteful as Pings only need to be sent if no other Packets have been sent in the
// time specified by the KeepAlive. More sophisticated timing logic could be added later.
func (c *Client) keepAlivePing() {
	ticker := time.NewTicker(time.Duration(c.Config.keepAlive) * time.Second)
	for {
		<-ticker.C
		err := c.SendPing()
		if err != nil {
			log.Error().
				Err(err).
				Str("source", "goqtt").
				Msg("could not send Ping in keepAlive")
		}
	}
}
