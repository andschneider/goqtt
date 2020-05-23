package goqtt

import (
	"bytes"
	"io"
	"log"
	"net"
)

func SendPacket(c net.Conn, packet []byte) {
	_, err := c.Write(packet)
	log.Printf("sent packet: %s", string(packet))
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
