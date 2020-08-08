package goqtt_test

import (
	"log"

	"github.com/andschneider/goqtt"
)

// This example subscribes to a MQTT broker. You will want to call ReadLoop afterwards
// to process incoming messages.
func ExampleClient_Subscribe() {
	// Create Client
	client := goqtt.NewClient("mqtt.eclipse.org")

	// Attempt a connection to the specified MQTT broker
	client.Connect()
	defer client.Disconnect()

	// Attempt to subscribe to the topic
	err := client.Subscribe()
	if err != nil {
		log.Fatalf("could not subscribe to topic: %v\n", err)
	}
}

// This example subscribes to a MQTT broker will print any incoming messages to std out.
func ExampleClient_SendPublish() {
	// Create Client
	client := goqtt.NewClient("mqtt.eclipse.org")

	// Attempt a connection to the specified MQTT broker
	client.Connect()
	defer client.Disconnect()

	// Attempt to publish a message
	err := client.SendPublish("hello world")
	if err != nil {
		log.Printf("could not send message: %v\n", err)
	}
}

// This example subscribes to a MQTT broker will print any incoming messages to std out, until the program exits.
func ExampleClient_ReadLoop() {
	// Create Client
	client := goqtt.NewClient("mqtt.eclipse.org")

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
