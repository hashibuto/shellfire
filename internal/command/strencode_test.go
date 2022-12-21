package command

import (
	"testing"

	"github.com/hashibuto/shellfire/internal/buffer"
	"github.com/hashibuto/shellfire/internal/utils"
)

func TestEncodeDecOut(t *testing.T) {
	StrEncodeCommand.Prepare()
	stdout, err := utils.CaptureStdout(func() error {
		return StrEncodeCommand.Process([]string{"-s=64", "/home/user/file.txt"})
	})
	if err != nil {
		t.Error(err)
		return
	}

	b := buffer.FromByteArray(stdout)
	expected := "0x2f686f6d652f7573    // \"/home/us\"\n0x65722f66696c652e    // \"er/file.\"\n0x7478740000000000    // \"txt     \"\n"
	if b.String() != expected {
		t.Errorf("Expected:\n%s\nDid not match actual\n%s\n", expected, b.String())
		return
	}
}
