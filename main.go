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
	conn, err := net.Dial("tcp", "192.168.1.189:1883")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	//create connection packet
	buf := new(bytes.Buffer)
	cpack := createConPacket()
	cpack.Write(buf)

	fmt.Println(buf.String())
	conn.Write(buf.Bytes())
	fmt.Println("connected")
	time.Sleep(1 * time.Second)

	//create subscription packet
	buf = new(bytes.Buffer)
	spack := createSubPacket()
	spack.Write(buf)
	fmt.Println(buf.String())
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
