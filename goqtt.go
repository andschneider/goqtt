package goqtt

import (
	"bytes"
	"log"
	"net"
)

func SendPacket(c net.Conn, packet []byte, verbose bool) {
	_, err := c.Write(packet)
	if verbose {
		log.Printf("sent packet: %s", string(packet))
	}
	if err != nil {
		log.Fatal(err)
	}
}

// TODO figure out how to generalize this to other packet types
func SendConnect(c net.Conn, verbose bool) error {
	cpack := CreateConnectPacket()
	buf := new(bytes.Buffer)
	err := cpack.Write(buf, verbose)
	if err != nil {
		log.Fatal(err)
		return err
	}
	SendPacket(c, buf.Bytes(), verbose)

	// response
	t, err := decodeByte(c)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if t != 32 {
		log.Printf("response type is not a connack packet. got %d", t)
	}

	cp := ConnackPacket{}
	err = cp.ReadConnackPacket(c)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if verbose {
		log.Printf("connack packet: %v", cp)
	}

	return nil
}

// TODO figure out how to generalize this to other packet types
func SendSubscribe(c net.Conn, topic string, verbose bool) error {
	spack := CreateSubscribePacket(topic)
	buf := new(bytes.Buffer)
	err := spack.Write(buf, verbose)
	if err != nil {
		log.Fatal(err)
		return err
	}
	SendPacket(c, buf.Bytes(), verbose)

	// response
	t, err := decodeByte(c)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if t != 144 {
		log.Printf("response type is not a suback packet. got %d", t)
	}

	sp := SubackPacket{}
	err = sp.ReadSubackPacket(c)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if verbose {
		log.Printf("suback packet: %v", sp)
	}

	return nil
}

func SubscribeLoop(conn net.Conn) {
	for {
		//log.Println("start loop")
		pp := PublishPacket{}
		p, err := pp.ReadPublishPacket(conn)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("TOPIC: %s MESSAGE: %s\n", p.Topic, string(p.Message))
	}
}
