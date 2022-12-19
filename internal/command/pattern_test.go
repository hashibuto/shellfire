package command

import (
	"testing"
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
