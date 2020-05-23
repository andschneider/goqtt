package goqtt

import (
	"bytes"
	"io"
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

// SendPing is a helper function to create a Ping and send it right away.
func SendPing(c net.Conn) {
	buf := new(bytes.Buffer)
	ping := CreatePingReqPacket()
	err := ping.Write(buf, true)
	if err != nil {
		log.Fatal(err)
	}
	SendPacket(c, buf.Bytes())
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
	// TODO this should be called in a timeout channel which handles ping responses as well
	//SendPing(conn)
	clear := make([]byte, 1)
	io.ReadFull(conn, clear)
	for {
		log.Println("start loop")
		//fh := make([]byte, 1)
		//_, err := io.ReadFull(conn, fh)

		//mtype := MessageTypesTemp[fh[0]]
		//log.Println("message type", mtype)

		//pp := PublishPacket{}
		//p, err := pp.ReadPublishPacket(conn)
		//pp := PingRespPacket{}
		//p, err := pp.ReadPingRespPacket(conn)
		//_, err = io.ReadFull(conn, []byte{40})

		_, err := Reader(conn)
		if err != nil {
			log.Fatal(err)
			return
		}
		//log.Printf("FIXED HEADER %v\n", fh)

		//log.Printf("%v\n", p.FixedHeader)

		//log.Printf("TOPIC: %s MESSAGE: %s\n", p.Topic, string(p.Message))
	}
}
