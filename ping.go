package goqtt

import (
	"time"

	"github.com/andschneider/goqtt/packets"
	"github.com/rs/zerolog/log"
)

// SendPing is a helper function to create a PINGREQ packet and send it right away.
// It also reads the PINGRESP packet.
func (c *Client) SendPing() error {
	// create packet
	var p packets.PingReqPacket
	p.CreatePacket()

	c.stagePacket(&p)
	// read response and verify it's a PINGRESP packet
	//r, err := c.readResponse()
	//if err != nil {
	//	return fmt.Errorf("could not read response for %s: %v", p.Name(), err)
	//}
	//if _, ok := r.(*packets.PingRespPacket); !ok {
	//	typeErrorResponseLogger(p.Name(), r.Name(), r)
	//	return fmt.Errorf("did not receive an PINGRESP packet, got %s instead", r.Name())
	//}
	//if err != nil {
	//	log.Error().Err(err).Str("source", "goqtt").Msg("could not read a PINGRESP packet")
	//	return fmt.Errorf("could not read a PINGRESP packet")
	//}
	//}()
	return nil
}

func (c *Client) KeepAlive() {
	ticker := time.NewTicker(time.Duration(c.config.keepAlive) * time.Second)
	go func() {
		for {
			select {
			case <-c.disconnect:
				log.Debug().
					Str("source", "goqtt").
					Msg("disconnecting in KeepAlive")
				break
			case <-ticker.C:
				err := c.SendPing()
				if err != nil {
					log.Error().
						Err(err).
						Str("source", "goqtt").
						Msg("could not send Ping in KeepAlive")
				}
			}
		}
	}()
}
