package goqtt

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/andschneider/goqtt/packets"
)

// TODO figure out how to generalize the Send<Packet> functions

// SendPublish sends a given message to a given topic in a PUBLISH packet.
func SendPublish(c net.Conn, topic string, message string, verbose bool) error {
	// create packet
	ppack := packets.CreatePublishPacket(topic, message)
	err := ppack.Write(c)
	if verbose {
		fmt.Printf("publish string: %v\n", &ppack)
	}
	if err != nil {
		return fmt.Errorf("could not write PUBLISH packet: %v", err)
	}

	// has no response with QOS 0
	return nil
}

// SendPing is a helper function to create a PINGREQ packet and send it right away.
// It also reads the PINGRESP packet.
func SendPing(c net.Conn, verbose bool) error {
	// create packet
	ping := packets.CreatePingReqPacket()
	err := ping.Write(c)
	if verbose {
		fmt.Printf("ping string: %v\n", &ping)
	}
	if err != nil {
		return fmt.Errorf("could not write PING packet: %v", err)
	}

	// response
	// why did i make this a goroutine?
	go func() {
		_, err = packets.Reader(c)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("ping resp %v", p)
	}()
	return nil
}

// SendConnect sends a CONNECT packet and reads the CONNACK response.
func SendConnect(c net.Conn, verbose bool) error {
	// create packet
	cpack := packets.CreateConnectPacket()
	err := cpack.Write(c)
	if verbose {
		fmt.Printf("connect string: %v\n", &cpack)
	}
	if err != nil {
		return fmt.Errorf("could not write CONNECT packet: %v", err)
	}

	// response
	_, err = packets.Reader(c)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// SendSubscribe sends a SUBSCRIBE packet to a given topic and reads the SUBACK packet.
func SendSubscribe(c net.Conn, topic string, verbose bool) error {
	// create packet
	spack := packets.CreateSubscribePacket(topic)
	err := spack.Write(c)
	if verbose {
		fmt.Printf("subscribe string: %v\n", &spack)
	}
	if err != nil {
		return fmt.Errorf("could not write SUBSCRIBE packet: %v", err)
	}

	// response
	_, err = packets.Reader(c)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// SubscribeLoop keeps a connection alive after a successful subscription to a topic and reads any incoming messages.
// It sends pings every 30 seconds to keep the connection alive.
func SubscribeLoop(conn net.Conn) {
	ticker := time.NewTicker(30 * time.Second)
	// TODO add disconnect functionality
	disconnect := make(chan bool)
	go func() {
		for {
			select {
			case <-disconnect:
				return
			case <-ticker.C:
				err := SendPing(conn, false)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	for {
		//log.Println("start loop")
		// TODO add callback function to process packet from Reader
		_, err := packets.Reader(conn)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}
