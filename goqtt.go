package goqtt

import (
	"bytes"
	"flag"
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
	cpack := CreateConnectPacket()
	err = cpack.Write(buf, *verbose)
	if err != nil {
		log.Fatal(err)
	}

	sendPacket(conn, buf.Bytes())
	log.Printf("connected to %s:%s\n", *ip, *port)
	time.Sleep(1 * time.Second)

	// create publish packet
	buf.Reset()
	ppack := CreatePublishPacket(*topic, "hihihihihihi")
	err = ppack.Write(buf, *verbose)
	if err != nil {
		log.Fatal(err)
	}
	sendPacket(conn, buf.Bytes())

	// create subscription packet
	buf.Reset()
	spack := CreateSubscribePacket(*topic)
	err = spack.Write(buf, *verbose)
	if err != nil {
		log.Fatal(err)
	}
	sendPacket(conn, buf.Bytes())
	log.Printf("subscribed to %s\n", *topic)

	subscribeLoop(conn)
}

func sendPacket(c net.Conn, packet []byte) {
	_, err := c.Write(packet)
	log.Printf("sent packet: %s", string(packet))
	if err != nil {
		log.Fatal(err)
	}
}

func subscribeLoop(conn net.Conn) {
	for {
		//log.Println("start loop")
		pp := PublishPacket{}
		p, err := pp.ReadPublishPacket(conn)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("TOPIC: %s MESSAGE: %s\n", p.Topic, string(p.Message))
	}

}
