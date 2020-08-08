/*
This example subscribes to a MQTT broker will print any incoming messages to the terminal.

to run: go run ./examples/simple/sub/main.go

It is meant as a simple example showing the shortest path to subscribing to a topic. Things like
error checking/handling are purposefully left out. To see a longer example, check out the
'subscribe' example.
*/

package main

import (
	"log"

	"github.com/andschneider/goqtt"
)

func main() {
	// Create Client
	cfg := goqtt.ClientConfig{
		ClientId:  "simple-sub",
		KeepAlive: 30,
		Server:    "mqtt.eclipse.org",
		Port:      "1883",
		Topic:     "goqtt",
	}
	client, _ := goqtt.NewClient(&cfg)

	// Attempt a connection to the specified MQTT broker
	client.Connect()
	defer client.Disconnect()

	// Attempt to subscribe to the topic
	client.Subscribe()

	// Read messages indefinitely
	for {
		log.Println("waiting for message")
		m, _ := client.ReadLoop()
		if m != nil {
			log.Printf("received message: '%s' from topic: '%s'", string(m.Message), m.Topic)
		}
	}
}
