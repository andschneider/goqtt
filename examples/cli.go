package main

import (
	"flag"
	"os"
)

// CLI helper used in multiple examples
func cli() (*string, *string, *string, *bool) {
	ip := flag.String("ip", "", "IP address to connect to.")
	port := flag.String("port", "", "Port of host.")
	topic := flag.String("topic", "", "Topic(s) to subscribe to.")
	verbose := flag.Bool("v", false, "Verbose output. Default is false.")
	flag.Parse()

	if *ip == "" || *port == "" || *topic == "" {
		flag.Usage()
		os.Exit(1)
	}
	return ip, port, topic, verbose
}
