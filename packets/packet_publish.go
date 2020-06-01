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

func (p *PublishPacket) String() string {
	return fmt.Sprintf("%v topic: %s messageid: %v message: %s", p.FixedHeader, p.Topic, p.MessageId, p.Message)
}

func CreatePublishPacket(topic string, message string) (pp PublishPacket) {
	pp.FixedHeader = FixedHeader{PacketType: publishType}
	pp.Topic = topic
	pp.MessageId = []byte{0, 1}
	pp.Message = []byte(message)
	return
}

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

func (p *PublishPacket) ReadPublishPacket(r io.Reader) (*PublishPacket, error) {
	var fh FixedHeader
	fh.PacketType = publishType
	err := fh.read(r)
	if err != nil {
		return nil, err
	}

	packetBytes := make([]byte, fh.RemainingLength)
	n, err := io.ReadFull(r, packetBytes)
	if err != nil {
		return nil, err
	}
	if n != fh.RemainingLength {
		return nil, fmt.Errorf("failed to read expected data from PUBLISH packet: %v", err)
	}

	mes := bytes.NewBuffer(packetBytes)
	p.FixedHeader = fh
	p.Topic, err = decodeString(mes)
	if err != nil {
		return nil, err
	}

	// would change if QOS wasn't always 0
	var payloadLength = p.RemainingLength
	payloadLength -= len(p.Topic) + 2
	if payloadLength < 0 {
		return nil, fmt.Errorf("error unpacking PUBLISH payload, payload length < 0")
	}
	p.Message = make([]byte, payloadLength)
	_, err = mes.Read(p.Message)
	return p, err
}
