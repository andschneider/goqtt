package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	// TODO: cli arguments
	ip := "192.168.1.189"
	port := "1883"
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	verbose := false // TODO: cli argument
	// create connection packet
	buf := new(bytes.Buffer)
	cpack := CreateConnectPacket()
	cpack.Write(buf, verbose)

	conn.Write(buf.Bytes())
	fmt.Println("connected")
	time.Sleep(1 * time.Second)

	// create subscription packet
	topic := "test/topic" // TODO: cli argument
	buf = new(bytes.Buffer)
	spack := CreateSubscribePacket(topic)
	spack.Write(buf, verbose)
	fmt.Printf("Subscribe response %s\n", buf.String())
	conn.Write(buf.Bytes())
	for {
		mustCopy(os.Stdout, conn)
	}
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
