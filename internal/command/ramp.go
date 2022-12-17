package command

import (
	"fmt"

	"github.com/hashibuto/artillery"
)

var RampCommand = &artillery.Command{
	Name:        "ramp",
	Description: "generates a buffer of characters for visual buffer recognition",
	Arguments: []*artillery.Argument{
		{
			Name:        "length",
			Description: "buffer length in bytes",
			Type:        artillery.Int,
		},
	},
	Options: []*artillery.Option{
		{
			Name:        "hex",
			ShortName:   'h',
			Description: "output in hex format for visualization purposes",
			Type:        artillery.Bool,
			Value:       true,
		},
		{
			Name:        "align",
			ShortName:   'a',
			Description: "align the ramp by prepending n bytes of data",
			Type:        artillery.Int,
		},
	},
	OnExecute: generateRamp,
}

// generateRamp generates a visualization pattern for easy recognition of the buffer on the stack
func generateRamp(n artillery.Namespace, p *artillery.Processor) error {
	var args struct {
		Align  int
		Hex    bool
		Length int
	}
	err := artillery.Reflect(n, &args)
	if err != nil {
		return err
	}

	if args.Align > args.Length {
		return fmt.Errorf("Aligment bytes must be fewer than the total buffer length")
	}

	size := args.Length + 4 + args.Align
	buffer := make([]byte, size)
	for i := 0; i < args.Align; i++ {
		buffer[i] = 0xFF
	}
	offset := 0
	for i := args.Align; i < len(buffer)-4; i += 4 {
		if offset%2 == 0 {
			buffer[i] = 1
			buffer[i+1] = 1
			buffer[i+2] = 1
			buffer[i+3] = 1
		} else {
			buffer[i] = byte(offset >> 24)
			buffer[i+1] = byte(offset >> 16)
			buffer[i+2] = byte(offset >> 8)
			buffer[i+3] = byte(offset)

			// Skip zeros
			for j := 0; j < 4; j++ {
				if buffer[i+j] == 0 {
					buffer[i+j] = 0xEE
				}
			}
		}

		offset++
	}

	for _, char := range buffer[:args.Length] {
		if args.Hex {
			fmt.Printf("\\x%02x", char)
		} else {
			fmt.Printf("%c", char)
		}
	}

	if args.Hex {
		fmt.Println()
	}

	return nil
}
