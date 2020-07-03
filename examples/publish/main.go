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
	"net"
	"os"

	"github.com/andschneider/goqtt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	server := flag.String("server", "mqtt.eclipse.org", "Server to connect to.")
	port := flag.String("port", "1883", "Port of host.")
	topic := flag.String("topic", "hello/world", "Topic(s) to subscribe to.")
	message := flag.String("message", "hello", "Message to send to topic.")
	verbose := flag.Bool("v", false, "Verbose output. Default is false.")
	flag.Parse()

	// Set logger to pretty print instead of structured json
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Set log level to debug if verbose is passed in
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	conn, err := net.Dial("tcp", *server+":"+*port)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer conn.Close()

	err = goqtt.SendConnect(conn)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Str("server", *server).Str("port", *port).Msg("connection successful")

	// create publish packet
	err = goqtt.SendPublish(conn, *topic, *message)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Str("TOPIC", *topic).Str("MESSAGE", *message).Msg("message sent successfully")
}
