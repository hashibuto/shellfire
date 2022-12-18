package command

import (
	"testing"
)

func TestUniqueness(t *testing.T) {
	seen := map[string]struct{}{}
	seq := generateByteSeq(max)
	for i := 0; i < max-4; i++ {
		str := string(seq[i : i+4])
		_, exists := seen[str]
		if exists {
			t.Errorf("Produced a non-unique sequence at offset %d", i)
			return
		}
		seen[str] = struct{}{}
	}
}
