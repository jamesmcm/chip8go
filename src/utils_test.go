package main

import "testing"
import "fmt"

func TestROMFromString(t *testing.T) {
	// ADD Vx
	var romarray [4]byte
	rombytes := ROMFromString("a2cc 6a06")

	copy(romarray[:], rombytes)
	if romarray != [4]byte{0xA2, 0xCC, 0x6A, 0x06} {
		t.Errorf("ROMFromString incorrect, got: %s, want: %s",
			fmt.Sprint(romarray), fmt.Sprint([4]byte{0xA2, 0xCC, 0x6A, 0x06}))
	}
}
