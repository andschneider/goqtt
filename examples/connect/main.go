/*
This example simply connects to a MQTT broker and then terminates.

to run: go run ./examples/connect/main.go

The default broker is the publicly available server hosted by the Eclipse foundation, but can be changed by specifying a
different host name or IP address with the -server flag.
*/

package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/andschneider/goqtt"
)

func main() {
	server := flag.String("server", "mqtt.eclipse.org", "Server to connect to.")
	port := flag.String("port", "1883", "Port of host.")
	topic := flag.String("topic", "hello/world", "Topic to subscribe/unsubscribe to.")
	verbose := flag.Bool("v", false, "Verbose output. Default is false.")
	flag.Parse()

	conn, err := net.Dial("tcp", *server+":"+*port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	err = goqtt.SendConnect(conn, *verbose)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("connected to %s:%s\n", *server, *port)

	err = goqtt.SendSubscribe(conn, *topic, *verbose)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("subscribed to %s\n", *topic)

	sleep := 3 * time.Second
	log.Printf("sleeping for %s\n", sleep)
	time.Sleep(sleep)

	err = goqtt.SendUnsubscribe(conn, *topic, *verbose)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("unsubscribed from %s\n", *topic)

	sleep = 3 * time.Second
	log.Printf("sleeping for %s\n", sleep)
	time.Sleep(sleep)

	log.Println("sending a disconnect request")
	err = goqtt.SendDisconnect(conn, *verbose)
	if err != nil {
		log.Fatal(err)
		return
	}
}
