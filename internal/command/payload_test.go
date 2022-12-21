package command

import (
	"testing"

	"github.com/hashibuto/shellfire/internal/buffer"
	"github.com/hashibuto/shellfire/internal/utils"
)

func TestBuildPayload(t *testing.T) {
	PayloadCommand.Prepare()
	stdout, err := utils.CaptureStdout(func() error {
		return PayloadCommand.Process([]string{"10", "\\x33\\x34\\x35\\x36", "\\x01\\x02\\x03\\x04\\x05\\x06"})
	})
	if err != nil {
		t.Error(err)
		return
	}

	b := buffer.FromByteArray(stdout)
	actual := b.HexString()
	expected := "\\x90\\x90\\x01\\x02\\x03\\x04\\x05\\x06\\x90\\x90\\x36\\x35\\x34\\x33"
	if actual != expected {
		t.Errorf("Expected:\n%s\nDid not match actual\n%s\n", expected, actual)
		return
	}
}

func TestBuildPayloadUnevenBytes(t *testing.T) {
	PayloadCommand.Prepare()
	stdout, err := utils.CaptureStdout(func() error {
		return PayloadCommand.Process([]string{"11", "\\x33\\x34\\x35\\x36", "\\x01\\x02\\x03\\x04\\x05\\x06"})
	})
	if err != nil {
		t.Error(err)
		return
	}

	b := buffer.FromByteArray(stdout)
	actual := b.HexString()
	expected := "\\x90\\x90\\x01\\x02\\x03\\x04\\x05\\x06\\x90\\x90\\x90\\x36\\x35\\x34\\x33"
	if actual != expected {
		t.Errorf("Expected:\n%s\nDid not match actual\n%s\n", expected, actual)
		return
	}
}
