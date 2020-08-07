package packets

import (
	"bytes"
	"reflect"
	"testing"
)

const (
	testTopic   = "test/topic"
	testMessage = "hello world"
)

func TestFixedHeader(t *testing.T) {
	var expected bytes.Buffer
	expected.Write([]byte{16, 0})

	fh := FixedHeader{PacketType: connectType}
	h := fh.WriteHeader()
	if !bytes.Equal(expected.Bytes(), h.Bytes()) {
		t.Errorf("Expected %b, got %b\n", expected, h)
	}
}

var (
	testConnackPacket    = []byte{connackType.packetId, 2, 0, 0}
	testConnectPacket    = []byte{connectType.packetId, 2, 0, 0}
	testDisconnectPacket = []byte{disconnectType.packetId, 0}
	testPingReqPacket    = []byte{pingReqType.packetId, 0}
	testPingRespPacket   = []byte{pingRespType.packetId, 0}
	testPublishPacket    = []byte{publishType.packetId, 0}
	testSubackPacket     = []byte{subackType.packetId, 3, 0, 1, 0}
	testSubscribePacket  = []byte{subscribeType.packetId, 3, 0, 1, 0}
	testUnsubackPacket   = []byte{unsubackType.packetId, 2, 0, 1}
	//testUnsubscribePacket = []byte{unsubscribeType.packetId, 14, 0, 1, 0, 10, 116, 101, 115, 116, 47, 116, 111, 112, 105, 99}
	testUnsubscribePacket = append([]byte{unsubscribeType.packetId, 14, 0, 1, 0}, []byte("\ntest/topic")...)
)

// defaultPackets are what the packet's CreatePacket is expected to create.
var defaultPackets = map[string]Packet{
	"connack": &ConnackPacket{
		FixedHeader: FixedHeader{connackType, 2}},
	"connect": &ConnectPacket{},
	"disconnect": &DisconnectPacket{
		FixedHeader: FixedHeader{disconnectType, 0}},
	"pingreq": &PingReqPacket{
		FixedHeader: FixedHeader{pingReqType, 0}},
	"pingresp": &PingRespPacket{
		FixedHeader: FixedHeader{pingRespType, 0}},
	"publish": &PublishPacket{},
	"suback": &SubackPacket{
		FixedHeader: FixedHeader{subackType, 3},
		MessageId:   []byte{0, 1}, ReturnCodes: []byte{0}},
	"subscribe": &SubscribePacket{},
	"unsuback": &UnsubackPacket{
		FixedHeader: FixedHeader{unsubackType, 2},
		MessageId:   []byte{0, 1}},
	"unsubscribe": &UnsubscribePacket{
		FixedHeader: FixedHeader{unsubscribeType, 14},
		MessageId:   []byte{0, 1}, Topics: []string{testTopic}},
}

var testCases = []struct {
	packetType    PacketType
	packetBytes   []byte
	blankPacket   Packet
	defaultPacket Packet
}{
	{connackType, testConnackPacket, &ConnackPacket{}, defaultPackets["connack"]},
	{connectType, testConnectPacket, &ConnectPacket{}, defaultPackets["connect"]},
	{disconnectType, testDisconnectPacket, &DisconnectPacket{}, defaultPackets["disconnect"]},
	{pingReqType, testPingReqPacket, &PingReqPacket{}, defaultPackets["pingreq"]},
	{pingRespType, testPingRespPacket, &PingRespPacket{}, defaultPackets["pingresp"]},
	{publishType, testPublishPacket, &PublishPacket{}, defaultPackets["publish"]},
	{subackType, testSubackPacket, &SubackPacket{}, defaultPackets["suback"]},
	{subscribeType, testSubscribePacket, &SubscribePacket{}, defaultPackets["subscribe"]},
	{unsubscribeType, testUnsubscribePacket, &UnsubscribePacket{}, defaultPackets["unsubscribe"]},
	{unsubackType, testUnsubackPacket, &UnsubackPacket{}, defaultPackets["unsuback"]},
}

// Test the NewPacket function which should create a blank packet based on a packetId.
func TestNewPacket(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.packetType.name, func(t *testing.T) {
			got, err := NewPacket(tc.packetType.packetId)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tc.blankPacket, got) {
				t.Fatalf("expected %s, got %v", tc.blankPacket, got)
			}
		})
	}
}

// TODO add connect, publish, and subscribe packets
// Test the default values set by the ReadPacket methods.
func TestReadPacket(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.packetType.name, func(t *testing.T) {
			r := bytes.NewReader(tc.packetBytes)
			got, err := ReadPacket(r)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tc.defaultPacket, got) {
				t.Fatalf("expected %s, got %v", tc.defaultPacket, got)
			}
		})
	}
}

// TODO connect, publish, and subscribe packets
// Test the default values set by CreatePacket methods.
func TestCreatePacket(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.packetType.name, func(t *testing.T) {
			var buf bytes.Buffer
			var err error

			// call unique Create wrapper methods
			if tc.packetType == unsubscribeType {
				var up UnsubscribePacket
				up.CreateUnsubscribePacket(testTopic)
				err = up.Write(&buf)
			} else if tc.packetType == subscribeType {
				var sp SubscribePacket
				sp.CreateSubscribePacket(testTopic)
				err = sp.Write(&buf)
			} else if tc.packetType == publishType {
				var pp PublishPacket
				pp.CreatePublishPacket(testTopic, testMessage)
				err = pp.Write(&buf)
			} else {
				// call general CreatePacket methods
				tc.blankPacket.CreatePacket()
				err = tc.blankPacket.Write(&buf)
			}
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(tc.packetBytes, buf.Bytes()) {
				t.Fatalf("expected %s, got %v", tc.packetBytes, buf)
			}
		})
	}
}
