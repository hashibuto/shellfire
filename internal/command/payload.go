package command

import (
	"fmt"

	"github.com/hashibuto/artillery"
	"github.com/hashibuto/shellfire/internal/buffer"
	"github.com/hashibuto/shellfire/internal/utils"
)

var PayloadCommand = &artillery.Command{
	Name:        "payload",
	Description: "generates a shellcode payload for a buffer overflow",
	Arguments: []*artillery.Argument{
		{
			Name:        "offset",
			Description: "buffer offset where return address will be written",
			Type:        artillery.Int,
		},
		{
			Name:        "returnAddress",
			Description: "return address in highest-to-lowest hex format (eg. \\xed\\x32\\x44\\x55)",
		},
		{
			Name:        "shellcode",
			Description: "shellcode in hex string format (eg. \\xed\\x32\\x44\\x55\\xed\\x32\\x44\\x55)",
		},
	},
	Options: []*artillery.Option{
		{
			Name:        "append",
			ShortName:   'a',
			Description: "append shellcode after return address",
			Type:        artillery.Bool,
			Value:       true,
		},
		{
			Name:        "nsb",
			ShortName:   'n',
			Description: "number of nopsled bytes to insert (defaults to half the remaining buffer space)",
			Type:        artillery.Int,
		},
		{
			Name:        "hex",
			ShortName:   'h',
			Description: "output in hex format for visualization purposes",
			Type:        artillery.Bool,
			Value:       true,
		},
		{
			Name:        "big",
			ShortName:   'b',
			Description: "specify big-endian byte order",
			Type:        artillery.Bool,
			Value:       true,
		},
	},
	OnExecute: makePayload,
}

// makePayload generates a buffer payload containingthe provided shellcode
func makePayload(n artillery.Namespace, p *artillery.Processor) error {
	var args struct {
		Append        bool
		Big           bool
		Hex           bool
		Nsb           int
		Offset        int
		ReturnAddress string
		Shellcode     string
	}
	artillery.Reflect(n, &args)

	shellcodeBytes, err := utils.ParseHexArrayString(args.Shellcode)
	if err != nil {
		return err
	}
	returnAddressBytes, err := utils.ParseHexArrayString(args.ReturnAddress)
	if err != nil {
		return err
	}

	if len(shellcodeBytes) > args.Offset {
		return fmt.Errorf("Shellcode is larger than the target buffer")
	}

	if args.Nsb > args.Offset-len(shellcodeBytes) {
		return fmt.Errorf("Nopsled is too long for the provided buffer")
	}

	endOfReturn := args.Offset + len(returnAddressBytes)
	totalSize := endOfReturn
	if args.Append {
		totalSize += len(shellcodeBytes)
	}

	// Init with NOP
	b := buffer.NewBuffer(totalSize, 0x90)

	var offset int
	if args.Append {
		offset = endOfReturn
	} else {
		if args.Nsb > 0 {
			offset = args.Nsb
		} else {
			offset = int((args.Offset - len(shellcodeBytes)) / 2)
		}
	}

	b.CopyTo(shellcodeBytes, offset)
	if args.Big {
		b.CopyTo(returnAddressBytes, args.Offset)
	} else {
		b.RevCopyTo(returnAddressBytes, args.Offset)
	}

	err = b.Validate()
	if err != nil {
		return err
	}
	b.Stdout(args.Hex)
	return nil
}
