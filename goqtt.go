package goqtt

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/andschneider/goqtt/packets"
)

// TODO figure out how to generalize the Send<Packet> functions

// SendPublish sends a given message to a given topic in a PUBLISH packet.
func SendPublish(c net.Conn, topic string, message string, verbose bool) error {
	// create packet
	var pp packets.PublishPacket
	pp.CreatePacket(topic, message)
	err := pp.Write(c)
	if verbose {
		fmt.Printf("publish string: %v\n", pp.String())
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
	var p packets.PingReqPacket
	p.CreatePacket()
	err := p.Write(c)
	if verbose {
		fmt.Printf("ping string: %v\n", p.String())
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
	var cp packets.ConnectPacket
	cp.CreatePacket()
	err := cp.Write(c)
	if verbose {
		fmt.Printf("connect string: %v\n", cp.String())
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

// SendSubscribe sends a SUBSCRIBE packet for a given topic and reads the SUBACK packet.
func SendSubscribe(c net.Conn, topic string, verbose bool) error {
	// create packet
	var sp packets.SubscribePacket
	sp.CreatePacket(topic)
	err := sp.Write(c)
	if verbose {
		fmt.Printf("subscribe string: %v\n", sp.String())
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

// SendUnsubscribe sends an UNSUBSCRIBE packet for a given topic and reads the UNSUBACK packet.
func SendUnsubscribe(c net.Conn, topic string, verbose bool) error {
	// create packet
	var up packets.UnsubscribePacket
	up.CreatePacket(topic)
	err := up.Write(c)
	if verbose {
		fmt.Printf("unsubscribe string: %v\n", up.String())
	}
	if err != nil {
		return fmt.Errorf("could not write UNSUBSCRIBE packet: %v", err)
	}

	// response
	_, err = packets.Reader(c)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// SendDisconnect sends a DISCONNECT packet.
func SendDisconnect(c net.Conn, verbose bool) error {
	// create packet
	var dp packets.DisconnectPacket
	dp.CreatePacket()
	err := dp.Write(c)
	if verbose {
		fmt.Printf("disconnect string: %v\n", dp.String())
	}
	if err != nil {
		return fmt.Errorf("could not write DISCONNECT packet: %v", err)
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
		// TODO add callback function to process packet from Reader
		_, err := packets.Reader(conn)
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
