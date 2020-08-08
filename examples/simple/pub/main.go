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
	cfg := &goqtt.ClientConfig{
		ClientId:  "simple-pub",
		KeepAlive: 30,
		Server:    "mqtt.eclipse.org",
		Port:      "1883",
		Topic:     "goqtt",
	}
	client, _ := goqtt.NewClient(cfg)

	// Attempt a connection to the specified MQTT broker
	client.Connect()
	defer client.Disconnect()

	// Attempt to publish a message
	err := client.SendPublish("hello world")
	if err != nil {
		fmt.Printf("could not send message: %v\n", err)
	}
}
