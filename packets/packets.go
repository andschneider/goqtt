package packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const MQTT3 = 4 // 0000100 in binary

type Packet interface {
	//CreatePacket()
	String() string
	Write(io.Writer) error
	Read(io.Reader) error
}

// PacketType represents the human readable name and the byte representation of the first byte of the packet header
type PacketType struct {
	Name     string
	packetId byte
}

// FixedHeader represents the packet type and the length of the remaining payload in the packet
type FixedHeader struct {
	PacketType
	RemainingLength int
}

func (fh *FixedHeader) String() string {
	return fmt.Sprintf("%s remaining length: %d", fh.PacketType.Name, fh.RemainingLength)
}

func (fh *FixedHeader) WriteHeader() (header bytes.Buffer) {
	header.WriteByte(fh.PacketType.packetId)
	header.Write(encodeLength(fh.RemainingLength))
	return
}

// Assumes the message type byte has already been read
func (fh *FixedHeader) read(r io.Reader) (err error) {
	fh.RemainingLength, err = decodeLength(r)
	return err
}

// ReadPacket tries to read a Packet from a TCP connection.
func ReadPacket(r io.Reader) (Packet, error) {
	packetId, err := decodeByte(r)
	if err != nil {
		return nil, err
	}

	p, err := NewPacket(packetId)
	if err != nil {
		return nil, err
	}

	err = p.Read(r)
	return p, err
}

// NewPacket creates an empty Packet according to the packetId parameter. The FixedHeader is set later in the packet's
// Read method
func NewPacket(packetId byte) (Packet, error) {
	switch packetId {
	case connackType.packetId:
		return &ConnackPacket{}, nil
	case connectType.packetId:
		return &ConnectPacket{}, nil
	case disconnectType.packetId:
		return &DisconnectPacket{}, nil
	case pingReqType.packetId:
		return &PingReqPacket{}, nil
	case pingRespType.packetId:
		return &PingRespPacket{}, nil
	case publishType.packetId:
		return &PublishPacket{}, nil
	case subackType.packetId:
		return &SubackPacket{}, nil
	case subscribeType.packetId:
		return &SubscribePacket{}, nil
	case unsubackType.packetId:
		return &UnsubackPacket{}, nil
	case unsubscribeType.packetId:
		return &UnsubscribePacket{}, nil
	}
	return nil, fmt.Errorf("packet type not accounted for: %v\n", packetId)
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

func decodeByte(b io.Reader) (byte, error) {
	num := make([]byte, 1)
	_, err := b.Read(num)
	if err == io.EOF {
		// TODO do i care?
		//fmt.Println("EOF but i dont care")
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("could not read bytes in decodeByte: %v", err)
	}
	return num[0], nil
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

func decodeMessageId(b io.Reader) ([]byte, error) {
	var bb []byte
	num := make([]byte, 2)
	_, err := b.Read(num)
	if err != nil {
		return nil, err
	}
	return append(bb, num...), nil
}

func decodeUint16(b io.Reader) (uint16, error) {
	num := make([]byte, 2)
	_, err := b.Read(num)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(num), nil
}
