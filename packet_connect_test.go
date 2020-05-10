package main

import (
	"fmt"
	"os"
	"testing"
)

func TestConnectPacket(t *testing.T) {
	cp := CreateConnectPacket()
	cp.Write(os.Stdout, true)
	fmt.Println()
}
