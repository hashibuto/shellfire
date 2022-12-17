package command

import "testing"

func TestParseHexString(t *testing.T) {
	input := "\\x03\\xF1\\Xfa"
	output, err := parseHexString(input)
	if err != nil {
		t.Error(err)
		return
	}

	if len(output) != 3 {
		t.Errorf("Expected an output length of 3 bytes")
		return
	}

	if output[0] != 3 || output[1] != 241 || output[2] != 250 {
		t.Errorf("Output did not match input")
		return
	}
}
