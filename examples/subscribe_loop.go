package main

import (
	"bytes"
	"github.com/andschneider/goqtt"
	"log"
	"net"
	"time"
)

func main() {
	ip, port, topic, verbose := cli()

	conn, err := net.Dial("tcp", *ip+":"+*port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// create connection packet
	buf := new(bytes.Buffer)
	cpack := goqtt.CreateConnectPacket()
	err = cpack.Write(buf, *verbose)
	if err != nil {
		log.Fatal(err)
	}

	goqtt.SendPacket(conn, buf.Bytes())
	log.Printf("connected to %s:%s\n", *ip, *port)
	time.Sleep(1 * time.Second)

	// create publish packet
	buf.Reset()
	ppack := goqtt.CreatePublishPacket(*topic, "hihihihihihi")
	err = ppack.Write(buf, *verbose)
	if err != nil {
		log.Fatal(err)
	}
	goqtt.SendPacket(conn, buf.Bytes())

	// create subscription packet
	buf.Reset()
	spack := goqtt.CreateSubscribePacket(*topic)
	err = spack.Write(buf, *verbose)
	if err != nil {
		log.Fatal(err)
	}
	goqtt.SendPacket(conn, buf.Bytes())
	log.Printf("subscribed to %s\n", *topic)

	goqtt.SubscribeLoop(conn)
}
