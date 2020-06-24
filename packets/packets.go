package packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

const MQTT3 = 4 // 0000100 in binary

// PacketType represents the human readable name and the byte representation of the first byte of the packet header
type PacketType struct {
	name     string
	packetId byte
}

// FixedHeader represents the packet type and the length of the remaining payload in the packet
type FixedHeader struct {
	PacketType
	RemainingLength int
}

func (fh *FixedHeader) String() string {
	return fmt.Sprintf("%s remaining length: %d", fh.PacketType.name, fh.RemainingLength)
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

// Reader reads in a packet from a TCP connection, determining packet type based on the first byte of the packet.
func Reader(r io.Reader) (interface{}, error) {
	pid, err := decodeByte(r)
	if err != nil {
		return nil, fmt.Errorf("could not decode byte from fixed header while reading packet: %v", err)
	}
	switch pid {
	case connectType.packetId:
		//fmt.Println("GOT A CONNECT")
		cp := ConnectPacket{}
		err = cp.Read(r)
		if err != nil {
			return nil, fmt.Errorf("could not read %s packet: %v", connectType.name, err)
		}
		//fmt.Printf("connect packet: %v\n", &cp)
		return cp, nil
	case connackType.packetId:
		//fmt.Println("GOT A CONNACK RESPONSE")
		cp := ConnackPacket{}
		err = cp.Read(r)
		if err != nil {
			return nil, fmt.Errorf("could not read CONNACK packet: %v", err)
		}
		//log.Printf("connack packet: %v", cp)
	case pingReqType.packetId:
		//fmt.Println("GOT A PING REQUEST")
		pr := PingReqPacket{}
		err := pr.Read(r)
		if err != nil {
			return nil, fmt.Errorf("could not read %s packet: %v", pingReqType.name, err)
		}
		return pr, nil
	case pingRespType.packetId:
		//log.Println("GOT A PING RESPONSE")
		pr := PingRespPacket{}
		err := pr.Read(r)
		if err != nil {
			return nil, fmt.Errorf("could not read PINGRESP packet: %v", err)
		}
	case publishType.packetId:
		//fmt.Println("GOT A PUBLISH RESPONSE")
		pp := PublishPacket{}
		p, err := pp.Read(r)
		if err != nil {
			return nil, fmt.Errorf("could not read PUBLISH packet: %v", err)
		}
		// TODO replace this with a callback function
		log.Printf("TOPIC: %s MESSAGE: %s\n", p.Topic, string(p.Message))
		return pp, nil
	case subscribeType.packetId:
		sp := SubscribePacket{}
		err := sp.Read(r)
		if err != nil {
			return nil, fmt.Errorf("could not read %s packet: %v", subscribeType.name, err)
		}
		return sp, nil
	case subackType.packetId:
		//fmt.Println("GOT A SUBACK RESPONSE")
		sp := SubackPacket{}
		err = sp.Read(r)
		if err != nil {
			return nil, fmt.Errorf("could not read SUBACK packet: %v", err)
		}
		//log.Printf("suback packet: %v\n", &sp)
	case unsubscribeType.packetId:
		//fmt.Println("GOT AN UNSUBSCRIBE REQUEST")
		up := UnsubscribePacket{}
		err = up.Read(r)
		if err != nil {
			return nil, fmt.Errorf("could not read UNSUBSCRIBE packet: %v", err)
		}
		//log.Printf("unsubscribe packet: %s\n", up.String())
		return up, nil
	case unsubackType.packetId:
		//fmt.Println("GOT AN UNSUBACK RESPONSE")
		up := UnsubackPacket{}
		err = up.Read(r)
		if err != nil {
			return nil, fmt.Errorf("could not read UNSUBACK packet: %v", err)
		}
		//log.Printf("unsuback packet: %s\n", up.String())
		//return up, nil
	case disconnectType.packetId:
		//fmt.Println("got a disconnect")
		return DisconnectPacket{}, nil
	case 0:
		return nil, io.EOF
	default:
		log.Printf("Type read in was not accounted for: %v\n", pid)
	}

	return nil, nil
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
