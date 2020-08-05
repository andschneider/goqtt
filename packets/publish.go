package packets

import (
	"bytes"
	"fmt"
	"io"
)

type PublishPacket struct {
	FixedHeader
	Topic     string
	MessageId []byte
	Message   []byte
}

var publishType = PacketType{
	name:     "PUBLISH",
	packetId: 48,
}

// Name returns the packet type name.
func (p *PublishPacket) Name() string {
	return p.name
}

// CreatePacket creates a new packet with the appropriate FixedHeader.
// It sets default values where needed as well.
func (p *PublishPacket) CreatePacket(topic string, message string) {
	p.FixedHeader = FixedHeader{PacketType: publishType}
	p.Topic = topic
	p.MessageId = []byte{0, 1}
	p.Message = []byte(message)
}

func (p *PublishPacket) String() string {
	return fmt.Sprintf("%v topic: %s messageid: %v message: %s", p.FixedHeader, p.Topic, p.MessageId, p.Message)
}

// Write creates the bytes.Buffer of the packet and writes them to
// the supplied io.Writer.
func (p *PublishPacket) Write(w io.Writer) error {
	var body bytes.Buffer
	var err error

	body.Write(encodeString(p.Topic))
	//body.Write(p.MessageId)  // Packet ID is only set on QOS 1 or 2, which I'm not using right now

	p.RemainingLength = body.Len() + len(p.Message)
	packet := p.WriteHeader()
	packet.Write(body.Bytes())
	packet.Write(p.Message)
	_, err = packet.WriteTo(w)

	return err
}

// Read creates the packet from an io.Reader. It assumes that the
// first byte, the packet id, has already been read.
func (p *PublishPacket) Read(r io.Reader) error {
	var fh FixedHeader
	fh.PacketType = publishType
	err := fh.read(r)
	if err != nil {
		return err
	}

	packetBytes := make([]byte, fh.RemainingLength)
	n, err := io.ReadFull(r, packetBytes)
	if err != nil {
		return err
	}
	if n != fh.RemainingLength {
		return fmt.Errorf("failed to read expected data from PUBLISH packet: %v", err)
	}

	mes := bytes.NewBuffer(packetBytes)
	p.FixedHeader = fh
	p.Topic, err = decodeString(mes)
	if err != nil {
		return err
	}

	// would change if QOS wasn't always 0
	var payloadLength = p.RemainingLength
	payloadLength -= len(p.Topic) + 2
	if payloadLength < 0 {
		return fmt.Errorf("error unpacking PUBLISH payload, payload length < 0")
	}
	p.Message = make([]byte, payloadLength)
	_, err = mes.Read(p.Message)
	if err != nil {
		return err
	}
	return nil
}
