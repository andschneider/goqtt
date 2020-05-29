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
	ProtocolName    string //[]byte
	ProtocolVersion byte
	ConnectFlags    byte
	KeepAlive       []byte

	ClientIdentifier string
}

var connectType = PacketType{
	name:     "CONNECT",
	packetId: 16,
}

func (c *ConnectPacket) String() string {
	return fmt.Sprintf("%v protocolname: %v protocolversion: %v connectflags: %08b clientid: %s", c.FixedHeader, c.ProtocolName, c.ProtocolVersion, c.ConnectFlags, c.ClientIdentifier)
}

func CreateConnectPacket() (cp ConnectPacket) {
	cp.FixedHeader = FixedHeader{PacketType: connectType}
	cp.ProtocolName = "MQTT"
	cp.ProtocolVersion = MQTT3
	cp.ConnectFlags = 2
	cp.KeepAlive = []byte{0, 60}
	hostname, _ := os.Hostname()
	cp.ClientIdentifier = hostname + strconv.Itoa(time.Now().Second())
	return
}

func (c *ConnectPacket) Write(w io.Writer) error {
	var body bytes.Buffer
	var err error

	body.Write(encodeString(c.ProtocolName))
	body.WriteByte(c.ProtocolVersion)
	body.WriteByte(c.ConnectFlags)
	body.Write(c.KeepAlive)
	body.Write(encodeString(c.ClientIdentifier))

	c.RemainingLength = body.Len()
	packet := c.WriteHeader()
	packet.Write(body.Bytes())
	_, err = packet.WriteTo(w)

	return err
}

func (c *ConnectPacket) Read(r io.Reader) error {
	var fh FixedHeader
	fh.PacketType = connectType
	err := fh.read(r)
	if err != nil {
		return fmt.Errorf("could not read in header: %v", err)
	}
	c.FixedHeader = fh
	c.ProtocolName, err = decodeString(r)
	if err != nil {
		return err
	}
	c.ProtocolVersion, err = decodeByte(r)
	if err != nil {
		return err
	}
	c.ConnectFlags, err = decodeByte(r)
	if err != nil {
		return err
	}
	c.KeepAlive, err = decodeMessageId(r)
	if err != nil {
		return err
	}
	c.ClientIdentifier, err = decodeString(r)
	if err != nil {
		return err
	}
	return nil
}
