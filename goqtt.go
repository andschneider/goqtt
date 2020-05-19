package main

import (
	"bytes"
	"flag"
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
	cpack := CreateConnectPacket()
	cpack.Write(buf, *verbose)

	sendPacket(conn, buf.Bytes())
	log.Printf("connected to %s:%s\n", *ip, *port)
	time.Sleep(1 * time.Second)

	// create subscription packet
	buf.Reset()
	spack := CreateSubscribePacket(*topic)
	spack.Write(buf, *verbose)
	sendPacket(conn, buf.Bytes())
	log.Printf("subscribed to %s\n", *topic)
	//for {
	//	mustCopy(os.Stdout, conn)
	//	time.Sleep(5)
	//	fmt.Println()
	//}
	subscribeLoop(conn)
}

func sendPacket(c net.Conn, packet []byte) {
	_, err := c.Write(packet)
	if err != nil {
		log.Fatal(err)
	}
}

func subscribeLoop(conn net.Conn) {
	//r := bufio.NewReader(conn)
	for {
		log.Println("start loop")
		message := make([]byte, int(26))
		_, err := io.ReadFull(conn, message)
		//line, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("%s\n", message)
	}

}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
