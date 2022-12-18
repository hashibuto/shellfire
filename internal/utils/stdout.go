package utils

import (
	"fmt"
	"os"
)

// HexOut writes data in hex string format to the stdout
func HexOut(data []byte) {
	for i := 0; i < len(data); i++ {
		_, err := fmt.Printf("\\x%02x", data[i])
		if err != nil {
			panic(err)
		}
	}
}

// BinOut writes binary data to the stdout
func BinOut(data []byte) {
	_, err := os.Stdout.Write(data)
	if err != nil {
		panic(err)
	}
}

// Write sends either bin or hex string data to the stdout
func Write(data []byte, hex bool) {
	if hex {
		HexOut(data)
	} else {
		BinOut(data)
	}
}

// WriteValue writes a value of length to the stdout
func WriteValue(length int, value byte, hex bool) {
	if length == 0 {
		return
	}

	buffer := make([]byte, length)
	for i := 0; i < length; i++ {
		buffer[i] = value
	}

	Write(buffer, hex)
}
