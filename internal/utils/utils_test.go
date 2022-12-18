package utils

import "testing"

func TestParseHex(t *testing.T) {
	value, err := ParseNumber("\\x331F")
	if err != nil {
		t.Error(err)
		return
	}

	if value != 13087 {
		t.Errorf("Decimal value did not match expected")
	}
}

func TestParseHex2(t *testing.T) {
	value, err := ParseNumber("0x331F")
	if err != nil {
		t.Error(err)
		return
	}

	if value != 13087 {
		t.Errorf("Decimal value did not match expected")
	}
}

func TestParseDec(t *testing.T) {
	value, err := ParseNumber("2345")
	if err != nil {
		t.Error(err)
		return
	}

	if value != 2345 {
		t.Errorf("Decimal value did not match expected: %d", value)
	}
}
