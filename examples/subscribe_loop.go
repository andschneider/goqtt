package main

import (
	"github.com/andschneider/goqtt"
	"log"
	"net"
)

func main() {
	ip, port, topic, verbose := cli()

	conn, err := net.Dial("tcp", *ip+":"+*port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	err = goqtt.SendConnect(conn, *verbose)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("connected to %s:%s\n", *ip, *port)

	err = goqtt.SendSubscribe(conn, *topic, *verbose)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("subscribed to %s\n", *topic)

	goqtt.SubscribeLoop(conn)
}
