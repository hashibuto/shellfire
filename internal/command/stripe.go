package command

import (
	"fmt"

	"github.com/hashibuto/artillery"
)

var StripeCommand = &artillery.Command{
	Name:        "stripe",
	Description: "generates a buffer of characters for visual buffer recognition",
	Arguments: []*artillery.Argument{
		{
			Name:        "offset",
			Description: "offset of return address in bytes",
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
			Description: "align the stripe pattern by prepending n bytes of data",
			Type:        artillery.Int,
		},
	},
	OnExecute: generateStripe,
}

// generateStripe generates a visualization pattern for easy recognition of the buffer on the stack
func generateStripe(n artillery.Namespace, p *artillery.Processor) error {
	var args struct {
		Align  int
		Hex    bool
		Offset int
	}
	err := artillery.Reflect(n, &args)
	if err != nil {
		return err
	}

	if args.Align > args.Offset {
		return fmt.Errorf("Aligment bytes must be fewer than the total buffer length")
	}

	size := args.Offset + 4 + args.Align
	buffer := make([]byte, size)
	for i := 0; i < args.Align; i++ {
		buffer[i] = 0xFF
	}
	offset := 0
	for i := args.Align; i < len(buffer)-4; i += 4 {
		var char byte
		if offset%2 == 0 {
			char = 0x11
		} else {
			char = 0x22
		}
		for j := 0; j < 4; j++ {
			buffer[i+j] = char
		}

		offset++
	}

	for i := 0; i < 4; i++ {
		buffer[args.Offset+i] = 'A'
	}
	for _, char := range buffer[:args.Offset+8] {
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
