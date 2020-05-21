package goqtt

import (
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
