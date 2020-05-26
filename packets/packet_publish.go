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

func CreatePublishPacket(topic string, message string) (pp PublishPacket) {
	pp.FixedHeader = FixedHeader{MessageType: "PUBLISH"}
	pp.Topic = topic
	pp.MessageId = []byte{0, 1}
	pp.Message = encodeString(message)
	return
}

func (p *PublishPacket) Write(w io.Writer, v bool) error {
	var body bytes.Buffer
	var err error

	body.Write(encodeString(p.Topic))
	//body.Write(p.MessageId)  // Packet ID is only set on QOS 1 or 2, which I'm not using right now

	p.FixedHeader.RemainingLength = body.Len() + len(p.Message)
	packet := p.FixedHeader.WriteHeader()
	packet.Write(body.Bytes())
	packet.Write(p.Message)

	if v {
		fmt.Println("BODY", body)
		fmt.Println("PACKET", packet)
	}
	_, err = packet.WriteTo(w)

	return err
}

func (p *PublishPacket) ReadPublishPacket(r io.Reader) (*PublishPacket, error) {
	var fh FixedHeader
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
		return nil, fmt.Errorf("failed to read expected data")
	}

	mes := bytes.NewBuffer(packetBytes)
	p.FixedHeader = fh
	p.Topic, err = decodeString(mes)
	if err != nil {
		return nil, err
	}

	// would change if QOS wasn't always 0
	var payloadLength = p.FixedHeader.RemainingLength
	payloadLength -= len(p.Topic) + 2
	if payloadLength < 0 {
		return nil, fmt.Errorf("error unpacking publish, payload length < 0")
	}
	p.Message = make([]byte, payloadLength)
	_, err = mes.Read(p.Message)
	return p, err
}
