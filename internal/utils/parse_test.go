package utils

import "testing"

func TestNormalInput(t *testing.T) {
	input := "\\x02\\x03\\x04\\xFA\\Xfd"
	b, err := ParseHexArrayString(input)
	if err != nil {
		t.Error(err)
		return
	}
	expected := []byte{2, 3, 4, 0xfa, 0xfd}
	if len(b) != len(expected) {
		t.Errorf("Output does not match input length")
	}
	for i := 0; i < len(expected); i++ {
		if expected[i] != b[i] {
			t.Errorf("Output sequence\n%v\ndoes not match input sequence\n", b)
			return
		}
	}
}

func TestBadInput1(t *testing.T) {
	input := "\\x02 \\x03\\x04"
	_, err := ParseHexArrayString(input)
	if err == nil {
		t.Error("Parser didn't catch bad input")
		return
	}
}

func TestBadInput2(t *testing.T) {
	input := "\\x02\\x03\\x04/x34"
	_, err := ParseHexArrayString(input)
	if err == nil {
		t.Error("Parser didn't catch bad input")
		return
	}
}

func TestBadInput3(t *testing.T) {
	input := "\\x02\\x03\\x04//34"
	_, err := ParseHexArrayString(input)
	if err == nil {
		t.Error("Parser didn't catch bad input")
		return
	}
}

func TestBadInput4(t *testing.T) {
	input := "\\x02\\x03\\x044"
	_, err := ParseHexArrayString(input)
	if err == nil {
		t.Error("Parser didn't catch bad input")
		return
	}
}

func TestBadInput5(t *testing.T) {
	input := "\\x02\\x03\\x04G"
	_, err := ParseHexArrayString(input)
	if err == nil {
		t.Error("Parser didn't catch bad input")
		return
	}
}
