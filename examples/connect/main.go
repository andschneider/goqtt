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
	"github.com/andschneider/goqtt/packets"
)

func main() {
	server := flag.String("server", "mqtt.eclipse.org", "Server to connect to.")
	port := flag.String("port", "1883", "Port of host.")
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

	sleep := 3 * time.Second
	log.Printf("sleeping for %s\n", sleep)
	time.Sleep(sleep)

	log.Println("sending a disconnect request")
	dp := packets.CreateDisconnectPacket()
	err = dp.Write(conn)
	if err != nil {
		log.Fatal(err)
		return
	}
}
