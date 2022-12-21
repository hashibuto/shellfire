package buffer

import (
	"fmt"
	"os"
	"strings"
)

type Buffer struct {
	buffer []byte
	length int
}

// NewBuffer creates a new buffer object of the specified length
func NewBuffer(length int, init byte) *Buffer {
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = init
	}
	return &Buffer{
		buffer: b,
		length: length,
	}
}

// FromByteArray creates a new buffer from a byte array
func FromByteArray(b []byte) *Buffer {
	buf := NewBuffer(len(b), 0)
	buf.CopyTo(b, 0)
	return buf
}

// String returns a string representation of the buffer
func (buffer *Buffer) String() string {
	return string(buffer.buffer)
}

// Stripe will draw a stripe pattern in the target buffer of a length of length bytes starting from start.  The stripe length will
// be stripeLength which is half the diversity of the stripe.  If the stripe length is 4 it will look like 0000111100001111
func (buffer *Buffer) Stripe(start int, length int, stripeLength int) {
	offset := 0
	var char byte
	for i := start; i < start+length; i += 4 {
		if offset%2 == 0 {
			char = 'D'
		} else {
			char = 'U'
		}
		offset++

		for j := 0; j < stripeLength; j++ {
			if i+j >= buffer.length {
				// Prevent overdraw
				return
			}

			buffer.buffer[i+j] = char
		}
	}
}

// Paint will paint a character on the buffer from start position for length characters
func (buffer *Buffer) Paint(start int, length int, char byte) {
	for i := start; i < start+length; i++ {
		if i >= buffer.length {
			// Prevent overdraw
			return
		}
		buffer.buffer[i] = char
	}
}

// CopyTo copies as much of source as possible onto destination at the specified offset
func (buffer *Buffer) CopyTo(source []byte, destOffset int) {
	for i := 0; i < len(source); i++ {
		if i+destOffset >= buffer.length {
			return
		}
		buffer.buffer[i+destOffset] = source[i]
	}
}

// RevCopyTo copies as much of source as possible onto destination at the specified offset in reverse order
func (buffer *Buffer) RevCopyTo(source []byte, destOffset int) {
	sourceLen := len(source)
	for i := 0; i < len(source); i++ {
		if i+destOffset >= buffer.length {
			return
		}
		buffer.buffer[i+destOffset] = source[sourceLen-i-1]
	}
}

// HexString writes data in hex string format to the stdout
func (buffer *Buffer) HexString() string {
	strs := []string{}
	for i := 0; i < buffer.length; i++ {
		strs = append(strs, fmt.Sprintf("\\x%02x", buffer.buffer[i]))
	}

	return strings.Join(strs, "")
}

// SafeString returns a string with only printable characters
func (buffer *Buffer) SafeString() string {
	clone := buffer.buffer[:]
	for i := 0; i < len(clone); i++ {
		char := clone[i]
		if char < 33 || char > 126 {
			clone[i] = ' '
		}
	}
	return string(clone)
}

// Bytes writes binary data to the stdout
func (buffer *Buffer) Bytes() []byte {
	return buffer.buffer
}

// Stdout sends either bin or hex string data to the stdout
func (buffer *Buffer) Stdout(hex bool) {
	if hex {
		fmt.Printf("%s", buffer.HexString())
	} else {
		os.Stdout.Write(buffer.buffer)
	}
}

// WriteByteAt writes a byte at the specified position
func (buffer *Buffer) WriteByteAt(offset int, char byte) {
	if offset < buffer.length {
		buffer.buffer[offset] = char
	}
}

// SubStr returns a sub string within the buffer
func (buffer *Buffer) SubStr(offset int, length int) string {
	if offset >= buffer.length {
		return ""
	}
	if offset+length >= buffer.length {
		length = buffer.length - offset
	}
	return string(buffer.buffer[offset : offset+length])
}

// Validate returns an error if null bytes are present in the buffer
func (buffer *Buffer) Validate() error {
	for i := 0; i < buffer.length; i++ {
		if buffer.buffer[i] == 0 {
			return fmt.Errorf("Null bytes present in buffer")
		}
	}

	return nil
}
