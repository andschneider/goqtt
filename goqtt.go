package goqtt

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/andschneider/goqtt/packets"
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
func sendPacket(c net.Conn, p packets.Packet, verbose bool) error {
	err := p.Write(c)
	if verbose {
		fmt.Printf("packet string: %v\n", p.String())
	}
	if err != nil {
		return fmt.Errorf("could not send Packet to TCP connection: %v", err)
	}
	return nil
}

// SendPublish sends a given message to a given topic in a PUBLISH packet.
func SendPublish(c net.Conn, topic string, message string, verbose bool) error {
	// create packet
	var p packets.PublishPacket
	p.CreatePacket(topic, message)

	err := sendPacket(c, &p, verbose)
	if err != nil {
		return fmt.Errorf("could not send %s packet: %v", p.Name, err)
	}

	// has no response with QOS 0
	return nil
}

// SendPing is a helper function to create a PINGREQ packet and send it right away.
// It also reads the PINGRESP packet.
func SendPing(c net.Conn, verbose bool) error {
	// create packet
	var p packets.PingReqPacket
	p.CreatePacket()

	err := sendPacket(c, &p, verbose)
	if err != nil {
		return fmt.Errorf("could not send %s packet: %v", p.Name, err)
	}

	// response
	// why did i make this a goroutine?
	go func() {
		_, err := readResponse(c)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return nil
}

// SendConnect sends a CONNECT packet and reads the CONNACK response.
func SendConnect(c net.Conn, verbose bool) error {
	// create packet
	var p packets.ConnectPacket
	p.CreatePacket()

	err := sendPacket(c, &p, verbose)
	if err != nil {
		return fmt.Errorf("could not send %s packet: %v", p.Name, err)
	}

	// response
	r, err := readResponse(c)
	if rp, ok := r.(*packets.ConnackPacket); !ok {
		fmt.Printf("received packet %v", rp.String())
		return fmt.Errorf("did not receive a CONNACK packet, got %s instead", rp.Name)
	}
	return err
}

// SendSubscribe sends a SUBSCRIBE packet for a given topic and reads the SUBACK packet.
func SendSubscribe(c net.Conn, topic string, verbose bool) error {
	// create packet
	var p packets.SubscribePacket
	p.CreatePacket(topic)

	err := sendPacket(c, &p, verbose)
	if err != nil {
		return fmt.Errorf("could not send %s packet: %v", p.Name, err)
	}

	// response
	r, err := readResponse(c)
	if rp, ok := r.(*packets.SubackPacket); !ok {
		fmt.Printf("received packet %v", rp.String())
		return fmt.Errorf("did not receive a SUBACK packet, got %s instead", rp.Name)
	}
	return err
}

// SendUnsubscribe sends an UNSUBSCRIBE packet for a given topic and reads the UNSUBACK packet.
func SendUnsubscribe(c net.Conn, topic string, verbose bool) error {
	// create packet
	var p packets.UnsubscribePacket
	p.CreatePacket(topic)

	err := sendPacket(c, &p, verbose)
	if err != nil {
		return fmt.Errorf("could not send %s packet: %v", p.Name, err)
	}

	// response
	r, err := readResponse(c)
	if rp, ok := r.(*packets.UnsubackPacket); !ok {
		fmt.Printf("received packet %v", rp.String())
		return fmt.Errorf("did not receive an UNSUBACK packet, got %s instead", rp.Name)
	}
	return err
}

// SendDisconnect sends a DISCONNECT packet.
func SendDisconnect(c net.Conn, verbose bool) error {
	// create packet
	var p packets.DisconnectPacket
	p.CreatePacket()

	err := sendPacket(c, &p, verbose)
	if err != nil {
		return fmt.Errorf("could not send %s packet: %v", p.Name, err)
	}

	return nil
}

// SubscribeLoop keeps a connection alive after a successful subscription to a topic and reads any incoming messages.
// It sends pings every 30 seconds to keep the connection alive.
func SubscribeLoop(conn net.Conn, verbose bool) {
	ticker := time.NewTicker(27 * time.Second)
	// TODO add disconnect functionality
	disconnect := make(chan bool)
	go func() {
		for {
			select {
			case <-disconnect:
				return
			case <-ticker.C:
				err := SendPing(conn, verbose)
				if err != nil {
					log.Fatal(err)
				}
				//fmt.Println("would send a ping")
			}
		}
	}()

	for {
		//log.Println("start loop")
		p, err := packets.ReadPacket(conn)
		// process packets based on type
		switch packet := p.(type) {
		case *packets.PublishPacket:
			log.Printf("TOPIC: %s MESSAGE: %s\n", packet.Topic, string(packet.Message))
		}
		if err != nil {
			if err == io.EOF {
				log.Println("Looks like the server closed the connection...")
				break
			}
			log.Fatal("subscribe loop error\n", err)
			return
		}
	}
}
