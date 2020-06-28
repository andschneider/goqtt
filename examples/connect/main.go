/*
This example simply connects to a MQTT broker and then terminates.

to run: go run ./examples/connect/main.go

The default broker is the publicly available server hosted by the Eclipse foundation, but can be changed by specifying a
different host name or IP address with the -server flag.
*/

package main

import (
	"flag"
	"net"
	"os"
	"time"

	"github.com/andschneider/goqtt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	server := flag.String("server", "mqtt.eclipse.org", "Server to connect to.")
	port := flag.String("port", "1883", "Port of host.")
	topic := flag.String("topic", "hello/world", "Topic to subscribe/unsubscribe to.")
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

	err = goqtt.SendSubscribe(conn, *topic)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Str("topic", *topic).Msg("subscribe successful")

	sleep := 3 * time.Second
	log.Info().Msgf("sleeping for %s\n", sleep)
	time.Sleep(sleep)

	err = goqtt.SendUnsubscribe(conn, *topic)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Str("topic", *topic).Msg("unsubscribe successful")

	sleep = 3 * time.Second
	log.Info().Msgf("sleeping for %s\n", sleep)
	time.Sleep(sleep)

	log.Info().Msg("sending a disconnect request")
	err = goqtt.SendDisconnect(conn)
	if err != nil {
		log.Fatal().Err(err)
	}
}
