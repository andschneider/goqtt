package goqtt

import (
	"fmt"
	"os"
	"testing"
)

func TestSubscribePacket(t *testing.T) {
	sp := CreateSubscribePacket(testTopic)
	// fmt.Println(sp)
	sp.Write(os.Stdout, true)
	fmt.Println()
}
