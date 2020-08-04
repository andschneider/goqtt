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

var testConnackPacket = []byte{connackType.packetId, byte(remainingLength), byte(sessionPresent), byte(returnCode)}

func TestReadPacket(t *testing.T) {
	testCases := []struct {
		name   string
		packet []byte
		want   Packet
	}{
		{"connack", testConnackPacket, &ConnackPacket{
			FixedHeader: FixedHeader{connackType, remainingLength},
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
