/*
This example publishes a message to a given topic

to run: go run ./examples/publish/main.go

The default broker is the publicly available server hosted by the Eclipse foundation, but can be changed by specifying a
different host name or IP address with the -server flag.

To change the topic and message, use the -topic and -message flags, respectively.
*/

package main

import (
	"flag"
	"log"
	"net"

	"github.com/andschneider/goqtt"
)

func main() {
	server := flag.String("server", "mqtt.eclipse.org", "Server to connect to.")
	port := flag.String("port", "1883", "Port of host.")
	topic := flag.String("topic", "hello/world", "Topic(s) to subscribe to.")
	message := flag.String("message", "hello", "Message to send to topic.")
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

	// create publish packet
	err = goqtt.SendPublish(conn, *topic, *message, *verbose)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("sent message: '%s' to topic: %s\n", *message, *topic)
}
