package packets

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type ConnectPacket struct {
	FixedHeader
	ProtocolName    []byte
	ProtocolVersion byte
	ConnectFlags    byte
	KeepAlive       []byte

	ClientIdentifier string
}

func CreateConnectPacket() (cp ConnectPacket) {
	cp.FixedHeader = FixedHeader{MessageType: "CONNECT"}
	cp.ProtocolName = []byte{0, 4, 77, 81, 84, 84} // "04MQTT"
	cp.ProtocolVersion = MQTT3
	cp.ConnectFlags = 2
	cp.KeepAlive = []byte{0, 60}
	hostname, _ := os.Hostname()
	cp.ClientIdentifier = hostname + strconv.Itoa(time.Now().Second())
	return
}

func (c *ConnectPacket) Write(w io.Writer, v bool) error {
	var body bytes.Buffer
	var err error

	body.Write(c.ProtocolName)
	body.WriteByte(c.ProtocolVersion)
	body.WriteByte(c.ConnectFlags)
	body.Write(c.KeepAlive)
	body.Write(encodeString(c.ClientIdentifier))

	c.FixedHeader.RemainingLength = body.Len()
	packet := c.FixedHeader.WriteHeader()
	packet.Write(body.Bytes())

	if v {
		fmt.Println("BODY", body)
		fmt.Println("PACKET", packet)
	}
	_, err = packet.WriteTo(w)

	return err
}
