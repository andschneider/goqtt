package goqtt

import (
	"bytes"
	"log"
	"net"
)

func SendPacket(c net.Conn, packet []byte, verbose bool) {
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
	ping := CreatePingReqPacket()
	err := ping.Write(buf, verbose)
	if err != nil {
		log.Fatal(err)
		return err
	}
	SendPacket(c, buf.Bytes(), verbose)

	// response
	go func() {
		_, err = Reader(c)
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
	cpack := CreateConnectPacket()
	err := cpack.Write(buf, verbose)
	if err != nil {
		log.Fatal(err)
		return err
	}
	SendPacket(c, buf.Bytes(), verbose)

	// response
	_, err = Reader(c)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// TODO figure out how to generalize this to other packet types
func SendSubscribe(c net.Conn, topic string, verbose bool) error {
	buf := new(bytes.Buffer)
	spack := CreateSubscribePacket(topic)
	err := spack.Write(buf, verbose)
	if err != nil {
		log.Fatal(err)
		return err
	}
	SendPacket(c, buf.Bytes(), verbose)

	// response
	_, err = Reader(c)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// TODO add callback function to process packet from Reader
func SubscribeLoop(conn net.Conn) {
	for {
		//log.Println("start loop")
		_, err := Reader(conn)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}
