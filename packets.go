package goqtt

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

const (
	MQTT5 = 5 // 0000101 in binary
	MQTT3 = 4 // 0000100 in binary
)

var MessageTypes = map[string]byte{
	"CONNECT":   16,  // 00010000 in binary
	"SUBSCRIBE": 130, // 10000010 in binary
	"PUBLISH":   48,  // 00110000 in binary
	"PINGREQ":   192, // 11000000 in binary
	"PINGRESP":  208, // 11010000 in binary
}

var MessageTypesTemp = map[byte]string{
	16:  "CONNECT",
	130: "SUBSCRIBE",
	48:  "PUBLISH",
	192: "PINGREQ",
	208: "PINGRESP",
}

type Packet struct {
}

type FixedHeader struct {
	MessageType     string
	RemainingLength int
}

func (fh *FixedHeader) WriteHeader() (header bytes.Buffer) {
	t, ok := MessageTypes[fh.MessageType]
	if !ok {
		fmt.Println("wrong message type, must be ...") // TODO: make this better
	}
	header.WriteByte(t)
	header.Write(encodeLength(fh.RemainingLength))
	return
}

func (fh *FixedHeader) read(r io.Reader) (err error) {
	fh.MessageType = "PUBLISH" // TODO generalize to different types
	fh.RemainingLength, err = decodeLength(r)
	return err
}

func (fh *FixedHeader) String() string {
	return fmt.Sprintf("%s remaining length: %d", fh.MessageType, fh.RemainingLength)
}

func Reader(r io.Reader) (*Packet, error) {
	var fh FixedHeader
	b := make([]byte, 1)

	_, err := io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("SDFSDF 0 %d\n", b[0])
	fh.MessageType = MessageTypesTemp[b[0]]

	switch fh.MessageType {
	case "PINGRESP":
		fmt.Println("GOT A PING RESPONSE")
	case "PUBLISH":
		fmt.Println("GOT A PUBLISH RESPONSE")

		pp := PublishPacket{}
		p, err := pp.ReadPublishPacket(r)
		if err != nil {
			return nil, err
		}
		log.Printf("TOPIC: %s MESSAGE: %s\n", p.Topic, string(p.Message))

	}
	//fmt.Printf("TYPETYPE %s\n", fh.MessageType)
	//err = fh.read(r)
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Printf("FHFHFH %v\n", fh)
	//fmt.Printf("FHFHFH %s %d\n", fh.MessageType, fh.RemainingLength)

	//packetBytes := make([]byte, fh.RemainingLength)
	//_, err = io.ReadFull(r, packetBytes)
	//if err != nil {
	//	return nil, err
	//}
	fmt.Println("done reading...")
	return nil, err
}

func encodeLength(length int) []byte {
	var encLength []byte
	for {
		digit := byte(length % 128)
		length /= 128
		if length > 0 {
			digit |= 0x80
		}
		encLength = append(encLength, digit)
		if length == 0 {
			break
		}
	}
	return encLength
}

func decodeLength(r io.Reader) (int, error) {
	var rLength uint32
	var multiplier uint32
	b := make([]byte, 1)
	for multiplier < 27 { //fix: Infinite '(digit & 128) == 1' will cause the dead loop
		_, err := io.ReadFull(r, b)
		if err != nil {
			return 0, err
		}

		digit := b[0]
		rLength |= uint32(digit&127) << multiplier
		if (digit & 128) == 0 {
			break
		}
		multiplier += 7
	}
	return int(rLength), nil
}

func encodeString(s string) []byte {
	b := []byte(s)
	e := make([]byte, 2)
	binary.BigEndian.PutUint16(e, uint16(len(b)))
	return append(e, b...)
}

func decodeString(b io.Reader) (string, error) {
	buf, err := decodeBytes(b)
	return string(buf), err
}

func decodeBytes(b io.Reader) ([]byte, error) {
	fieldLength, err := decodeUint16(b)
	if err != nil {
		return nil, err
	}

	field := make([]byte, fieldLength)
	_, err = b.Read(field)
	if err != nil {
		return nil, err
	}

	return field, nil
}

func decodeUint16(b io.Reader) (uint16, error) {
	num := make([]byte, 2)
	_, err := b.Read(num)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(num), nil
}
