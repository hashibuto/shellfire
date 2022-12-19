package command

import (
	"testing"

	"github.com/hashibuto/shellfire/internal/buffer"
	"github.com/hashibuto/shellfire/internal/utils"
)

func TestEval1(t *testing.T) {
	EvalCommand.Prepare()
	stdout, err := utils.CaptureStdout(func() error {
		return EvalCommand.Process([]string{"123 + 456 - 333"})
	})
	if err != nil {
		t.Error(err)
		return
	}

	b := buffer.FromByteArray(stdout)
	expected := "\\x000000f6\n"
	if b.String() != expected {
		t.Errorf("Expected:\n%s\nDid not match actual\n%s\n", expected, b.String())
		return
	}
}

func TestEval2(t *testing.T) {
	EvalCommand.Prepare()
	stdout, err := utils.CaptureStdout(func() error {
		return EvalCommand.Process([]string{"0x123 + \\x33FA - \\x333"})
	})
	if err != nil {
		t.Error(err)
		return
	}

	b := buffer.FromByteArray(stdout)
	expected := "\\x000031ea\n"
	if b.String() != expected {
		t.Errorf("Expected:\n%s\nDid not match actual\n%s\n", expected, b.String())
		return
	}
}

func TestEvalDecOut(t *testing.T) {
	EvalCommand.Prepare()
	stdout, err := utils.CaptureStdout(func() error {
		return EvalCommand.Process([]string{"-d", "0x123 + \\x33FA - \\x333"})
	})
	if err != nil {
		t.Error(err)
		return
	}

	b := buffer.FromByteArray(stdout)
	expected := "12778\n"
	if b.String() != expected {
		t.Errorf("Expected:\n%s\nDid not match actual\n%s\n", expected, b.String())
		return
	}
}
