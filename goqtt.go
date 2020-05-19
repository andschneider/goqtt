package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	// TODO: cli arguments
	ip := "192.168.1.189"
	port := "1883"
	verbose := false       // TODO: cli argument
	topic := "hello/world" // TODO: cli argument

	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// create connection packet
	buf := new(bytes.Buffer)
	cpack := CreateConnectPacket()
	cpack.Write(buf, verbose)

	sendPacket(conn, buf.Bytes())
	fmt.Printf("connected to %s:%s\n", ip, port)
	time.Sleep(1 * time.Second)

	// create subscription packet
	buf.Reset()
	spack := CreateSubscribePacket(topic)
	spack.Write(buf, verbose)
	sendPacket(conn, buf.Bytes())
	fmt.Printf("subscribed to %s\n", topic)
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
		log.Println("asdf")
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
