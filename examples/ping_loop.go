package main

import (
	"github.com/andschneider/goqtt"
	"log"
	"net"
	"time"
)

type ping struct {
	err        error
	pingPacket goqtt.PingRespPacket
}

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

	// TODO this should apart of the subscribe loop function
	ticker := time.NewTicker(30 * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				err = goqtt.SendPing(conn, *verbose)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	// subscribe to topic and read messages
	goqtt.SubscribeLoop(conn)
	//ticker.Stop()
	//done <- true
}
