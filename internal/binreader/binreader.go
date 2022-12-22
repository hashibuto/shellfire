package binreader

import (
	"encoding/binary"
	"fmt"
)

var ReadPastEnd = fmt.Errorf("read past end of buffer")

// NewBinReader constructs a new binary reader
func NewBinReader(data []byte, order binary.ByteOrder) *BinReader {
	return &BinReader{
		data:   data,
		length: len(data),
		order:  order,
	}
}

type BinReader struct {
	length int
	data   []byte
	order  binary.ByteOrder
	offset int
}

// ReadByte returns a single byte from the stream
func (br *BinReader) ReadByte() (byte, error) {
	if br.offset+1 >= br.length {
		return 0, ReadPastEnd
	}

	b := br.data[br.offset]
	br.offset++
	return b, nil
}

// ReadWord returns a single word from the stream
func (br *BinReader) ReadWord() (uint16, error) {
	if br.offset+2 >= br.length {
		return 0, ReadPastEnd
	}

	var w uint16
	offset := br.offset
	if br.order == binary.LittleEndian {
		w = uint16(br.data[offset]) + (uint16(br.data[offset+1]) << 8)
	} else {
		w = uint16(br.data[offset+1]) + (uint16(br.data[offset]) << 8)
	}
	br.offset += 2
	return w, nil
}

// ReadDWord returns a single double-word from the stream
func (br *BinReader) ReadDWord() (uint32, error) {
	if br.offset+4 >= br.length {
		return 0, ReadPastEnd
	}

	var w uint32
	offset := br.offset
	if br.order == binary.LittleEndian {
		w = uint32(br.data[offset]) + (uint32(br.data[offset+1]) << 8) + (uint32(br.data[offset+2]) << 16) + (uint32(br.data[offset+3]) << 24)
	} else {
		w = uint32(br.data[offset+3]) + (uint32(br.data[offset+2]) << 8) + (uint32(br.data[offset+1]) << 16) + (uint32(br.data[offset]) << 24)
	}
	br.offset += 4
	return w, nil
}

// ReadQWord returns a single quad-word from the stream
func (br *BinReader) ReadQWord() (uint64, error) {
	if br.offset+8 >= br.length {
		return 0, ReadPastEnd
	}

	var w uint64
	offset := br.offset
	if br.order == binary.LittleEndian {
		w = uint64(br.data[offset]) + (uint64(br.data[offset+1]) << 8) + (uint64(br.data[offset+2]) << 16) + (uint64(br.data[offset+3]) << 24) + (uint64(br.data[offset+4]) << 32) + (uint64(br.data[offset+5]) << 40) + (uint64(br.data[offset+6]) << 48) + (uint64(br.data[offset+7]) << 56)
	} else {
		w = uint64(br.data[offset]+7) + (uint64(br.data[offset+6]) << 8) + (uint64(br.data[offset+5]) << 16) + (uint64(br.data[offset+4]) << 24) + (uint64(br.data[offset+3]) << 32) + (uint64(br.data[offset+2]) << 40) + (uint64(br.data[offset+1]) << 48) + (uint64(br.data[offset]) << 56)
	}
	br.offset += 8
	return w, nil
}

// ReadBytes reads and returns a slice of n bytes from the stream
func (br *BinReader) ReadBytes(num int) ([]byte, error) {
	if br.offset+num >= br.length {
		return nil, ReadPastEnd
	}
	b := br.data[br.offset : br.offset+num]
	br.offset += num
	return b, nil
}

// CopyBytes copies enough bytes to fill target, from the stream
func (br *BinReader) CopyBytes(target []byte) error {
	tLen := len(target)
	if br.offset+tLen >= br.length {
		return ReadPastEnd
	}
	offset := br.offset
	for i := 0; i < tLen; i++ {
		target[i] = br.data[offset+i]
	}
	br.offset += tLen

	return nil
}
