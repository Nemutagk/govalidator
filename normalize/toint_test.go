package normalize

import "testing"

func TestToInt_ParsesNumericString(t *testing.T) {
	got := ToInt(" 123 ", nil)
	if got != 123 {
		t.Fatalf("got %v, want 123", got)
	}
}

func TestToInt_ConvertsFloat64(t *testing.T) {
	got := ToInt(123.0, nil)
	if got != 123 {
		t.Fatalf("got %v, want 123", got)
	}
}

func TestToInt_NonNumericStringPassesThrough(t *testing.T) {
	got := ToInt("abc", nil)
	if got != "abc" {
		t.Fatalf("got %v, want unchanged 'abc'", got)
	}
}

func TestToInt_UnsupportedTypePassesThrough(t *testing.T) {
	got := ToInt(true, nil)
	if got != true {
		t.Fatalf("got %v, want true", got)
	}
}
