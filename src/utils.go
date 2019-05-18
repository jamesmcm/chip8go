package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func reverseMap(m map[uint16]uint8) map[uint8]uint16 {
	n := make(map[uint8]uint16)
	for k, v := range m {
		n[v] = k
	}
	return n
}

func charToHex(c rune) byte {
	switch {
	case c >= 48 && c <= 57:
		return byte(c - 48)
	case c >= 65 && c <= 70:
		return byte(c - 55)
	default:
		panic(fmt.Sprint("Bad Hex character: %i", c))
	}
}

// ROMFromString : Convert string of form "00 00 00" to a list of bytes
func ROMFromString(s string) []byte {
	var out []byte
	var curbyte byte
	var tempbyte byte

	s = strings.ToUpper(s)
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\n", "", -1)

	for i, byt := range s {
		tempbyte = charToHex(byt)

		if i%2 == 1 {
			curbyte <<= 4
			curbyte |= tempbyte
			out = append(out, curbyte)
		} else {
			curbyte = tempbyte
		}
	}
	return out
}

// PrintROM : Print list of bytes in "00 00" form
func PrintROM(rom []byte) {
	for i, byt := range rom {
		if i%2 == 0 {
			fmt.Printf("0x%03x: ", 0x200+i)
		}
		fmt.Printf("%02x", byt)
		if i%2 == 1 {
			fmt.Print("\n")
		}
	}
}

func readROM(filename string) []byte {
	dat, err := ioutil.ReadFile(filename)
	check(err)
	return dat
}
