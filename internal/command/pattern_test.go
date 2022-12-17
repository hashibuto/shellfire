package command

import (
	"fmt"
	"testing"
)

func TestUniqueness(t *testing.T) {
	seen := map[string]struct{}{}
	seq := generateByteSeq(max)
	fmt.Println(string(seq))
	for i := 0; i < max-4; i++ {
		str := string(seq[i : i+4])
		fmt.Println(str)
		_, exists := seen[str]
		if exists {
			t.Errorf("Produced a non-unique sequence at offset %d", i)
			return
		}
		seen[str] = struct{}{}
	}
}
