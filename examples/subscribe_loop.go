package main

import (
	"bytes"
	"flag"
	"github.com/andschneider/goqtt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	// CLI arguments
	ip := flag.String("ip", "", "IP address to connect to.")
	port := flag.String("port", "", "Port of host.")
	topic := flag.String("topic", "", "Topic(s) to subscribe to.")
	verbose := flag.Bool("v", false, "Verbose output. Default is false.")
	flag.Parse()

	if *ip == "" || *port == "" || *topic == "" {
		flag.Usage()
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", *ip+":"+*port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// create connection packet
	buf := new(bytes.Buffer)
	cpack := goqtt.CreateConnectPacket()
	err = cpack.Write(buf, *verbose)
	if err != nil {
		log.Fatal(err)
	}
	goqtt.SendPacket(conn, buf.Bytes())
	log.Printf("connected to %s:%s\n", *ip, *port)
	time.Sleep(1 * time.Second)

	// create subscription packet
	buf.Reset()
	spack := goqtt.CreateSubscribePacket(*topic)
	err = spack.Write(buf, *verbose)
	if err != nil {
		log.Fatal(err)
	}
	goqtt.SendPacket(conn, buf.Bytes())
	log.Printf("subscribed to %s\n", *topic)

	// subscribe to topic and read messages
	goqtt.SubscribeLoop(conn)
}
