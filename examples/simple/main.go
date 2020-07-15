/*
This example subscribes to a MQTT broker will print any incoming messages to the terminal.

to run: go run ./examples/simple/main.go

It is meant a simple example showing the shortest path to subscribing to a topic. Things like
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
	clientId := "simple"
	keepAlive := 30 // seconds
	broker := "mqtt.eclipse.org:1883"
	topic := "goqtt"

	client := goqtt.NewClient(goqtt.NewClientConfig(clientId, keepAlive, broker, topic))

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
