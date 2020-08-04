package packets

import (
	"bytes"
	"reflect"
	"testing"
)

const (
	testTopic = "test/topic"
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

func TestNewPacket(t *testing.T) {
	testCases := []struct {
		name     string
		packetId byte
		want     Packet
	}{
		{"connack", connackType.packetId, &ConnackPacket{}},
		{"connect", connectType.packetId, &ConnectPacket{}},
		{"disconnect", disconnectType.packetId, &DisconnectPacket{}},
		{"pingReq", pingReqType.packetId, &PingReqPacket{}},
		{"pingResp", pingRespType.packetId, &PingRespPacket{}},
		{"publish", publishType.packetId, &PublishPacket{}},
		{"suback", subackType.packetId, &SubackPacket{}},
		{"subscribe", subscribeType.packetId, &SubscribePacket{}},
		{"unsuback", unsubackType.packetId, &UnsubackPacket{}},
		{"unsubscribe", unsubscribeType.packetId, &UnsubscribePacket{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewPacket(tc.packetId)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected %s, got %v", tc.want, got)
			}
		})
	}
}

var (
	testConnackPacket    = []byte{connackType.packetId, byte(remainingLength), byte(sessionPresent), byte(returnCode)}
	testDisconnectPacket = []byte{disconnectType.packetId, 0}
	testPingReqPacket    = []byte{pingReqType.packetId, 0}
	testPingRespPacket   = []byte{pingRespType.packetId, 0}
	testSubackPacket     = []byte{subackType.packetId, 3, 0, 1, byte(returnCode)}
	testUnsubackPacket   = []byte{unsubackType.packetId, 2, 0, 1}
	//testUnsubscribePacket = []byte{unsubscribeType.packetId, 14, 0, 1, 0, 10, 116, 101, 115, 116, 47, 116, 111, 112, 105, 99}
	testUnsubscribePacket = append([]byte{unsubscribeType.packetId, 14, 0, 1, 0}, []byte("\ntest/topic")...)
)

func TestReadPacket(t *testing.T) {
	testCases := []struct {
		name   string
		packet []byte
		want   Packet
	}{
		{"connack", testConnackPacket, &ConnackPacket{
			FixedHeader: FixedHeader{connackType, remainingLength},
		}},
		{"disconnect", testDisconnectPacket, &DisconnectPacket{
			FixedHeader: FixedHeader{disconnectType, 0},
		}},
		{"pingreq", testPingReqPacket, &PingReqPacket{
			FixedHeader: FixedHeader{pingReqType, 0},
		}},
		{"pingresp", testPingRespPacket, &PingRespPacket{
			FixedHeader: FixedHeader{pingRespType, 0},
		}},
		{"suback", testSubackPacket, &SubackPacket{
			FixedHeader: FixedHeader{subackType, 3},
			MessageId:   []byte{0, 1}, ReturnCodes: []byte{0},
		}},
		{"unsuback", testUnsubackPacket, &UnsubackPacket{
			FixedHeader: FixedHeader{unsubackType, 2},
			MessageId:   []byte{0, 1},
		}},
		{"unsubscribe", testUnsubscribePacket, &UnsubscribePacket{
			FixedHeader: FixedHeader{unsubscribeType, 14},
			MessageId:   []byte{0, 1}, Topics: []string{testTopic},
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewReader(tc.packet)
			got, err := ReadPacket(r)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected %s, got %v", tc.want, got)
			}
		})
	}
}
