package goqtt

import (
	"bytes"
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
	err := ping.Write(buf, false)
	if err != nil {
		log.Fatal(err)
	}
	SendPacket(c, buf.Bytes())
}

func SubscribeLoop(conn net.Conn) {
	// TODO this should be called in a timeout channel which handles ping responses as well
	//SendPing(conn)
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
