package command

import "testing"

func TestEval(t *testing.T) {
	value, err := evalExpression(" 0x3328f2 + 123")
	if err != nil {
		t.Error(err)
		return
	}

	if value != 3352941 {
		t.Errorf("Expected diff value, got %d", value)
	}
}

func TestEval2(t *testing.T) {
	value, err := evalExpression("\\xFF - 1")
	if err != nil {
		t.Error(err)
		return
	}

	if value != 0xfe {
		t.Errorf("Expected diff value, got %d", value)
	}
}
