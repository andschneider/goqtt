/*
This example publishes a message to the given topic. If you don't have a subscriber running you
won't see anything happen.

to run: go run ./examples/simple/pub/main.go

It is meant as a simple example showing the shortest path to publishing a message. Things like
error checking/handling are purposefully left out. To see a longer example, check out the
'publish' example.
*/

package main

import (
	"fmt"

	"github.com/andschneider/goqtt"
)

func main() {
	// Create Client
	clientId := "simple-pub"
	keepAlive := 30 // seconds
	broker := "mqtt.eclipse.org:1883"
	topic := "goqtt"

	client := goqtt.NewClient(goqtt.NewClientConfig(clientId, keepAlive, broker, topic))

	// Attempt a connection to the specified MQTT broker
	client.Connect()
	defer client.Disconnect()

	// Attempt to publish a message
	err := client.SendPublish("hello world")
	if err != nil {
		fmt.Printf("could not send message: %v\n", err)
	}
}
