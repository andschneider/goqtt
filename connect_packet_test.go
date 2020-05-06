package main

import (
	"fmt"
	"os"
	"testing"
)

func TestConnectPacket(t *testing.T) {
	cp := create()
	// fmt.Println(cp)
	fo, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	cp.Write(fo)
	// cp.Write(os.Stdout)
	fmt.Println()
}
