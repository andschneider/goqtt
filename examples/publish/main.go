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
	"fmt"
	"os"

	"github.com/andschneider/goqtt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func config() (*goqtt.Client, string) {
	server := flag.String("server", "mqtt.eclipse.org", "Server to connect to.")
	port := flag.String("port", "1883", "Port of host.")
	topic := flag.String("topic", "goqtt", "Topic to subscribe/unsubscribe to.")
	message := flag.String("message", "hello", "Message to send to topic.")
	id := flag.String("id", "publish-example", "Client id. Default is 'publish-example' appended with the process id.")
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
	return client, *message
}

func main() {
	// parse command line args and create Client
	client, message := config()

	// Attempt a connection to the specified MQTT broker
	err := client.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to broker")
	}
	log.Info().Msg("connection successful")
	defer client.Disconnect()

	// create publish packet
	err = client.SendPublish(message)
	if err != nil {
		log.Fatal().Err(err).Msg("could not send message")
	}
	log.Info().Str("MESSAGE", message).Msg("message sent successfully")
}
