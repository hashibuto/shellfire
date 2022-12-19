package command

import (
	"testing"

	"github.com/hashibuto/shellfire/internal/buffer"
	"github.com/hashibuto/shellfire/internal/utils"
)

func TestUniqueness(t *testing.T) {
	seen := map[string]struct{}{}
	b := generateByteSeq(0, max, 0)
	for i := 0; i < max-4; i++ {
		str := b.SubStr(i, 4)
		_, exists := seen[str]
		if exists {
			t.Errorf("Produced a non-unique sequence at offset %d", i)
			return
		}
		seen[str] = struct{}{}
	}
}

func TestPattern1(t *testing.T) {
	PatternCommand.Prepare()
	stdout, err := utils.CaptureStdout(func() error {
		return PatternCommand.Process([]string{"create", "45"})
	})
	if err != nil {
		t.Error(err)
		return
	}

	b := buffer.FromByteArray(stdout)
	expected := "0000111122223333444455556666777788889999AAAAB"
	if b.String() != expected {
		t.Errorf("Expected:\n%s\nDid not match actual\n%s\n", expected, b.String())
		return
	}
}

func TestPattern2(t *testing.T) {
	PatternCommand.Prepare()
	stdout, err := utils.CaptureStdout(func() error {
		return PatternCommand.Process([]string{"create", "-f=10", "45"})
	})
	if err != nil {
		t.Error(err)
		return
	}

	b := buffer.FromByteArray(stdout)
	expected := "==========0000111122223333444455556666777788889999AAAAB"
	if b.String() != expected {
		t.Errorf("Expected:\n%s\nDid not match actual\n%s\n", expected, b.String())
		return
	}
}

func TestPattern3(t *testing.T) {
	PatternCommand.Prepare()
	stdout, err := utils.CaptureStdout(func() error {
		return PatternCommand.Process([]string{"create", "-f=10", "-c=Z", "45"})
	})
	if err != nil {
		t.Error(err)
		return
	}

	b := buffer.FromByteArray(stdout)
	expected := "ZZZZZZZZZZ0000111122223333444455556666777788889999AAAAB"
	if b.String() != expected {
		t.Errorf("Expected:\n%s\nDid not match actual\n%s\n", expected, b.String())
		return
	}
}
