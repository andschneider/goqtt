/*
This example subscribes to a MQTT broker will print any incoming messages to the terminal.

to run: go run ./examples/subscribe/main.go

The default broker is the publicly available server hosted by the Eclipse foundation, but can be changed by specifying a
different host name or IP address with the -server flag.

The default topic is "hello/world", which may or may not have any messages being published to it (if using the Eclipse
server). If nothing shows up, try a different topic or publish a message using the publish.go example.
*/

package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/andschneider/goqtt"
)

func main() {
	server := flag.String("server", "mqtt.eclipse.org", "Server to connect to.")
	port := flag.String("port", "1883", "Port of host.")
	topic := flag.String("topic", "hello/world", "Topic(s) to subscribe to.")
	verbose := flag.Bool("v", false, "Verbose output. Default is false.")
	flag.Parse()

	if *server == "" || *port == "" || *topic == "" {
		flag.Usage()
		os.Exit(1)
	}

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

	// subscribe to topic and read messages
	goqtt.SubscribeLoop(conn)
}
