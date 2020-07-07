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
	"fmt"
	"os"

	"github.com/andschneider/goqtt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func config() *goqtt.Client {
	server := flag.String("server", "mqtt.eclipse.org", "Server to connect to.")
	port := flag.String("port", "1883", "Port of host.")
	topic := flag.String("topic", "goqtt", "Topic to subscribe/unsubscribe to.")
	id := flag.String("id", "subscribe-example", "Client id. Default is 'subscribe-example' appended with the process id.")
	verbose := flag.Bool("v", false, "Verbose output. Default is false.")
	flag.Parse()

	// Set logger to pretty print instead of structured json
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Set log level to debug if verbose is passed in
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Create client
	clientId := fmt.Sprintf("%s-%d", *id, os.Getpid())
	keepAlive := 30 // seconds
	broker := fmt.Sprintf("%s:%s", *server, *port)
	log.Info().
		Str("clientId", clientId).
		Str("broker", broker).
		Int("keepAlive", keepAlive).
		Str("topic", *topic).
		Msg("client configuration")
	copts := goqtt.NewClientConfig(clientId, keepAlive, broker, *topic)
	client := goqtt.NewClient(copts)
	return client
}

func main() {
	// parse command line args and create Client
	client := config()

	// Attempt a connection to the specified MQTT broker
	err := client.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to broker")
	}
	log.Info().Msg("connection successful")

	// Attempt to subscribe to the topic
	err = client.Subscribe()
	if err != nil {
		log.Fatal().Err(err).Msg("could not subscribe to topic")
	}
	log.Info().Msg("subscription successful")

	// Read messages indefinitely
	//client.SubscribeLoop()
	client.KeepAlive()
	for {
		m, err := client.ReadLoop()
		if err != nil {
			log.Error().Err(err)
		}
		if m != nil {
			log.Info().
				Str("TOPIC", m.Topic).
				Str("DATA", string(m.Message)).
				Msg("publish packet received")
		}
	}
}
