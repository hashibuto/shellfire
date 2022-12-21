package command

import (
	"fmt"

	"github.com/hashibuto/artillery"
	"github.com/hashibuto/shellfire/internal/buffer"
)

var StrEncodeCommand = &artillery.Command{
	Name:        "strencode",
	Description: "encode a string as a series of hex values designe to fit into registers",
	Arguments: []*artillery.Argument{
		{
			Name:        "string",
			Description: "string to be encoded",
		},
	},
	Options: []*artillery.Option{
		{
			Name:        "size",
			ShortName:   's',
			Description: "size grouping in bits",
			Type:        artillery.Int,
			Default:     64,
		},
	},
	OnExecute: strEncodeCmd,
}

func strEncodeCmd(n artillery.Namespace, p *artillery.Processor) error {
	var args struct {
		String string
		Size   int
	}
	err := artillery.Reflect(n, &args)
	if err != nil {
		return err
	}

	if args.Size != 64 && args.Size != 32 && args.Size != 16 && args.Size != 8 {
		return fmt.Errorf("Size must be a multiple of 2 between 8 and 64")
	}

	bytes := args.Size / 8
	remainder := len(args.String) % bytes
	padding := bytes - remainder
	if padding > 0 {
		b := make([]byte, padding)
		for i := 0; i < padding; i++ {
			b[i] = 0
		}
		args.String += string(b)
	}

	for i := 0; i < len(args.String); i += bytes {
		fmt.Printf("0x")
		for j := 0; j < bytes; j++ {
			fmt.Printf("%02x", args.String[i+j])
		}
		b := buffer.FromByteArray([]byte(args.String[i : i+bytes]))
		fmt.Printf("    // \"%s\"\n", b.SafeString())
	}

	return nil
}
