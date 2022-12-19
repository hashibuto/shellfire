package buffer

import "testing"

func TestStripe(t *testing.T) {
	buffer := NewBuffer(30, '.')
	buffer.Stripe(5, 30, 4)

	expected := ".....DDDDUUUUDDDDUUUUDDDDUUUUD"
	if buffer.String() != expected {
		t.Errorf("Got %s expected %s", buffer.String(), expected)
		return
	}
}

func TestPaint(t *testing.T) {
	buffer := NewBuffer(30, '.')
	buffer.Paint(5, 30, 'B')

	expected := ".....BBBBBBBBBBBBBBBBBBBBBBBBB"
	if buffer.String() != expected {
		t.Errorf("Got %s expected %s", buffer.String(), expected)
		return
	}
}

func TestRevCopy(t *testing.T) {
	buffer := NewBuffer(5, 0)
	source := []byte("hello")
	buffer.RevCopyTo(source, 0)
	expected := "olleh"
	if buffer.String() != expected {
		t.Errorf("Got %s expected %s", buffer.String(), expected)
		return
	}
}

func TestSubStr(t *testing.T) {
	buffer := NewBuffer(5, 0)
	buffer.CopyTo([]byte("hello"), 0)
	expected := "llo"
	actual := buffer.SubStr(2, 10)
	if actual != expected {
		t.Errorf("Got %s expected %s", actual, expected)
		return
	}
}
