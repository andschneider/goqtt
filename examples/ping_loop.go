package main

import (
	"bytes"
	"flag"
	"github.com/andschneider/goqtt"
	"io"
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
	conack := make([]byte, 4)
	_, err = io.ReadFull(conn, conack)
	if err != nil {
		log.Fatal(err)
		return
	}

	// create subscription packet
	buf.Reset()
	spack := goqtt.CreateSubscribePacket(*topic)
	err = spack.Write(buf, *verbose)
	if err != nil {
		log.Fatal(err)
	}
	goqtt.SendPacket(conn, buf.Bytes())
	log.Printf("subscribed to %s\n", *topic)
	// suback
	suback := make([]byte, 4)
	_, err = io.ReadFull(conn, suback)
	if err != nil {
		log.Fatal(err)
		return
	}

	ticker := time.NewTicker(30 * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				goqtt.SendPing(conn)
				//log.Println("Tick at", t)
			}
		}
	}()
	//time.Sleep(1600 * time.Millisecond)
	//ticker.Stop()
	//done <- true
	// subscribe to topic and read messages
	goqtt.SubscribeLoop(conn)
	//for {
	//	mustCopy(os.Stdout, conn)
	//	//log.Printf("%s\n", message)
	//}
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
	//log.Printf()
}

