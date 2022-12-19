package command

import (
	"fmt"

	"github.com/hashibuto/artillery"
	"github.com/hashibuto/shellfire/internal/buffer"
	"github.com/hashibuto/shellfire/internal/utils"
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
		{
			Name:        "extend",
			ShortName:   'e',
			Description: "extend the stripe pattern after the return address by n bytes",
			Type:        artillery.Int,
		},
	},
	OnExecute: generateStripe,
}

// generateStripe generates a visualization pattern for easy recognition of the buffer on the stack
func generateStripe(n artillery.Namespace, p *artillery.Processor) error {
	var args struct {
		Align  int
		Extend int
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

	b := buffer.NewBuffer(args.Offset+utils.Arch32+args.Extend, 0)
	b.Paint(0, args.Align, '.')
	b.Stripe(args.Align, args.Offset-args.Align, utils.Arch32)
	b.Paint(args.Offset, utils.Arch32, 'w')
	b.Stripe(args.Offset+utils.Arch32, args.Extend, utils.Arch32)
	b.Stdout(args.Hex)

	return nil
}
