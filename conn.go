package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.1.189:1883")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	pack := create()

	buf := new(bytes.Buffer)
	pack.Write(buf)

	fmt.Println(buf.String())
	conn.Write(buf.Bytes())

	mustCopy(os.Stdout, conn)
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
