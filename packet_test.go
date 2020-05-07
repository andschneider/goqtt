package main

import (
	"fmt"
	"os"
	"testing"
)

func TestConnectPacket(t *testing.T) {
	cp := createConPacket()
	// fmt.Println(cp)
	fo, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	cp.Write(fo)
	// cp.Write(os.Stdout)
	fmt.Println()
}
func TestSubscribePacket(t *testing.T) {
	sp := createSubPacket()
	// fmt.Println(sp)
	sp.Write(os.Stdout)
	fmt.Println()
}
