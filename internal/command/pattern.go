package command

import (
	"fmt"

	"github.com/hashibuto/artillery"
	"github.com/hashibuto/shellfire/internal/buffer"
)

var baseChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWZYZabcdefghijklmnopqrstuvwxyz"
var baseLen = len(baseChars)
var max = 10000

var PatternCommand = &artillery.Command{
	Name:        "pattern",
	Description: "pattern generation / offset detection",
	SubCommands: []*artillery.Command{
		{
			Name:        "create",
			Description: "create an n-byte length pattern suitable for offset detection",
			Arguments: []*artillery.Argument{
				{
					Name:        "length",
					Description: "number of bytes to generate",
					Type:        artillery.Int,
				},
			},
			Options: []*artillery.Option{
				{
					Name:        "fixed",
					ShortName:   'f',
					Description: "fixed portion of pattern for extra large buffers",
					Type:        artillery.Int,
				},
				{
					Name:        "char",
					ShortName:   'c',
					Description: "character to use for fixed portion of pattern",
					Type:        artillery.String,
				},
			},
			OnExecute: generateBytes,
		},
		{
			Name:        "offset",
			Description: "report offset given the supplied pattern fragement",
			Arguments: []*artillery.Argument{
				{
					Name:        "pattern",
					Description: "pattern fragment from which to calculate the offset",
				},
			},
			Options: []*artillery.Option{
				{
					Name:        "fixed",
					ShortName:   'f',
					Description: "number of fixed bytes used in pattern creation",
					Type:        artillery.Int,
				},
				{
					Name:        "hex",
					ShortName:   'h',
					Description: "output offset in hex format",
					Type:        artillery.Bool,
					Value:       true,
				},
			},
			OnExecute: calculateOffset,
		},
	},
}

func NewIter() *Iter {
	return &Iter{
		byteSeq: []byte{},
		seen:    map[string]struct{}{},
	}
}

type Iter struct {
	seqNum  int
	seen    map[string]struct{}
	byteSeq []byte
}

func (iter *Iter) Next() byte {
	for {
		test := baseChars[iter.seqNum%baseLen]
		if len(iter.byteSeq) < 3 {
			iter.byteSeq = append(iter.byteSeq, test)
			return test
		}
		seq := iter.byteSeq[len(iter.byteSeq)-3 : len(iter.byteSeq)]
		seq = append(seq, test)
		testStr := string(seq)
		_, exists := iter.seen[testStr]
		if !exists {
			iter.byteSeq = append(iter.byteSeq, test)
			iter.seen[testStr] = struct{}{}
			return test
		}
		iter.seqNum++
	}
}

// generateBytes prints a character sequence to the stdout, representing a non-repeating pattern from which
// the true offset can be determined by inspecting a minimum of 4 consecutive bytes
func generateBytes(ns artillery.Namespace, processor *artillery.Processor) error {
	var args struct {
		Char   string
		Fixed  int
		Length int
	}
	err := artillery.Reflect(ns, &args)
	if err != nil {
		return err
	}

	if args.Length > max {
		return fmt.Errorf("Cannot produce patterns longer than %d bytes", max)
	}

	if args.Fixed > args.Length-4 {
		return fmt.Errorf("Fixed number of bytes must be reasonably smaller than the total pattern size")
	}

	if len(args.Char) > 1 {
		return fmt.Errorf("Fixed character must be exactly one character long")
	} else {
		args.Char = "="
	}

	b := generateByteSeq(args.Fixed, args.Length, args.Char[0])
	b.Stdout(false)
	return nil
}

// generateByteSeq generates a sequence of bytes of length in size
func generateByteSeq(fixedLength int, patternLength int, fillerByte byte) *buffer.Buffer {
	b := buffer.NewBuffer(fixedLength+patternLength, fillerByte)
	iter := NewIter()
	// Designed for minimum size of 32 bit registers (no problem with larger ie. 64 bit)
	for i := 0; i < patternLength; i++ {
		b.WriteByteAt(i+fixedLength, iter.Next())
	}

	return b
}

// calculateOffset determines the offset of the beginning of the supplied pattern, assuming it was generated by this program
func calculateOffset(ns artillery.Namespace, processor *artillery.Processor) error {
	var args struct {
		Hex     bool
		Fixed   int
		Pattern string
	}
	err := artillery.Reflect(ns, &args)
	if err != nil {
		return err
	}

	if len(args.Pattern) != 4 {
		return fmt.Errorf("Length of input string must be exactly 4 bytes")
	}
	pattern := args.Pattern
	iter := NewIter()
	seq := make([]byte, max)
	for i := 0; i < max; i++ {
		seq[i] = iter.Next()
		if i >= 4 {
			testStr := string(seq[i-4 : i])
			if pattern == testStr {
				foundAt := i - 4
				if args.Hex {
					fmt.Printf("\\x%08x\n", foundAt+args.Fixed)
				} else {
					fmt.Printf("%d\n", foundAt+args.Fixed)
				}
				return nil
			}
		}
	}

	return fmt.Errorf("Couldn't determine offset of provided pattern")
}
