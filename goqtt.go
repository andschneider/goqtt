package goqtt

import (
	"bytes"
	"github.com/andschneider/goqtt/packets"
	"log"
	"net"
	"time"
)

func sendPacket(c net.Conn, packet []byte, verbose bool) {
	_, err := c.Write(packet)
	if verbose {
		log.Printf("sent packet: %b", packet)
	}
	if err != nil {
		log.Fatal(err)
	}
}

// SendPing is a helper function to create a Ping and send it right away.
func SendPing(c net.Conn, verbose bool) error {
	buf := new(bytes.Buffer)
	ping := packets.CreatePingReqPacket()
	err := ping.Write(buf, verbose)
	if err != nil {
		log.Fatal(err)
		return err
	}
	sendPacket(c, buf.Bytes(), verbose)

	// response
	go func() {
		_, err = packets.Reader(c)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("ping resp %v", p)
	}()
	return nil
}

// TODO figure out how to generalize this to other packet types
func SendConnect(c net.Conn, verbose bool) error {
	buf := new(bytes.Buffer)
	cpack := packets.CreateConnectPacket()
	err := cpack.Write(buf, verbose)
	if err != nil {
		log.Fatal(err)
		return err
	}
	sendPacket(c, buf.Bytes(), verbose)

	// response
	_, err = packets.Reader(c)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// TODO figure out how to generalize this to other packet types
func SendSubscribe(c net.Conn, topic string, verbose bool) error {
	buf := new(bytes.Buffer)
	spack := packets.CreateSubscribePacket(topic)
	err := spack.Write(buf, verbose)
	if err != nil {
		log.Fatal(err)
		return err
	}
	sendPacket(c, buf.Bytes(), verbose)

	// response
	_, err = packets.Reader(c)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// TODO add callback function to process packet from Reader
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
		_, err := packets.Reader(conn)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}
