package goqtt

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/andschneider/goqtt/packets"
	"github.com/rs/zerolog/log"
)

// readResponse is a helper function to the ReadPacket function which attempts to read
// data from a TCP connection.
func readResponse(c net.Conn) (packets.Packet, error) {
	p, err := packets.ReadPacket(c)
	if err != nil {
		return nil, fmt.Errorf("could not read a Packet from the TCP connection: %v", err)
	}
	return p, nil
}

// sendPacket is a helper function to write a Packet to a TCP connection
func sendPacket(c net.Conn, p packets.Packet) error {
	log.Debug().
		Str("source", "goqtt").
		Str("packetType", p.Name()).
		Str("packet", p.String()).
		Msg("sending packet")
	err := p.Write(c)
	if err != nil {
		log.Error().
			Err(err).
			Str("source", "goqtt").
			Str("packetType", p.Name()).
			Str("packet", p.String()).
			Msg("could not send packet")
		return fmt.Errorf("could not send Packet to TCP connection: %v", err)
	}
	return nil
}

// SendPublish sends a given message to a given topic in a PUBLISH packet.
func SendPublish(c net.Conn, topic string, message string) error {
	// create packet
	var p packets.PublishPacket
	p.CreatePacket(topic, message)

	err := sendPacket(c, &p)
	if err != nil {
		return fmt.Errorf("could not send %s packet: %v", p.Name(), err)
	}

	// has no response with QOS 0
	return nil
}

// SendPing is a helper function to create a PINGREQ packet and send it right away.
// It also reads the PINGRESP packet.
func SendPing(c net.Conn) error {
	// create packet
	var p packets.PingReqPacket
	p.CreatePacket()

	err := sendPacket(c, &p)
	if err != nil {
		return fmt.Errorf("could not send %s packet: %v", p.Name(), err)
	}

	// response
	// why did i make this a goroutine?
	go func() {
		_, err := readResponse(c)
		if err != nil {
			log.Fatal().Err(err)
		}
	}()
	return nil
}

// SendConnect sends a CONNECT packet and reads the CONNACK response.
func SendConnect(c net.Conn) error {
	// create packet
	var p packets.ConnectPacket
	p.CreatePacket()

	err := sendPacket(c, &p)
	if err != nil {
		return fmt.Errorf("could not send %s packet: %v", p.Name(), err)
	}

	// response
	r, err := readResponse(c)
	if rp, ok := r.(*packets.ConnackPacket); !ok {
		typeErrorResponseLogger(p.Name(), rp.Name(), rp)
		return fmt.Errorf("did not receive a CONNACK packet, got %s instead", rp.Name())
	}
	return err
}

// SendSubscribe sends a SUBSCRIBE packet for a given topic and reads the SUBACK packet.
func SendSubscribe(c net.Conn, topic string) error {
	// create packet
	var p packets.SubscribePacket
	p.CreatePacket(topic)

	err := sendPacket(c, &p)
	if err != nil {
		return fmt.Errorf("could not send %s packet: %v", p.Name(), err)
	}

	// response
	r, err := readResponse(c)
	if rp, ok := r.(*packets.SubackPacket); !ok {
		typeErrorResponseLogger(p.Name(), rp.Name(), rp)
		return fmt.Errorf("did not receive a SUBACK packet, got %s instead", rp.Name())
	}
	return err
}

// SendUnsubscribe sends an UNSUBSCRIBE packet for a given topic and reads the UNSUBACK packet.
func SendUnsubscribe(c net.Conn, topic string) error {
	// create packet
	var p packets.UnsubscribePacket
	p.CreatePacket(topic)

	err := sendPacket(c, &p)
	if err != nil {
		return fmt.Errorf("could not send %s packet: %v", p.Name(), err)
	}

	// response
	r, err := readResponse(c)
	if rp, ok := r.(*packets.UnsubackPacket); !ok {
		typeErrorResponseLogger(p.Name(), rp.Name(), rp)
		return fmt.Errorf("did not receive an UNSUBACK packet, got %s instead", rp.Name())
	}
	return err
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

// SendDisconnect sends a DISCONNECT packet.
func SendDisconnect(c net.Conn) error {
	// create packet
	var p packets.DisconnectPacket
	p.CreatePacket()

	err := sendPacket(c, &p)
	if err != nil {
		return fmt.Errorf("could not send %s packet: %v", p.Name(), err)
	}

	return nil
}

// SubscribeLoop keeps a connection alive after a successful subscription to a topic and reads any incoming messages.
// It sends pings every 30 seconds to keep the connection alive.
func SubscribeLoop(conn net.Conn) {
	ticker := time.NewTicker(27 * time.Second)
	// TODO add disconnect functionality
	disconnect := make(chan bool)
	go func() {
		for {
			select {
			case <-disconnect:
				return
			case <-ticker.C:
				err := SendPing(conn)
				if err != nil {
					log.Fatal().Err(err)
				}
				//fmt.Println("would send a ping")
			}
		}
	}()

	for {
		p, err := packets.ReadPacket(conn)
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
