package normalize

import "testing"

func TestToFloat_ParsesNumericString(t *testing.T) {
	got := ToFloat(" 12.5 ", nil)
	if got != 12.5 {
		t.Fatalf("got %v, want 12.5", got)
	}
}

func TestToFloat_ConvertsInt(t *testing.T) {
	got := ToFloat(12, nil)
	if got != 12.0 {
		t.Fatalf("got %v, want 12.0", got)
	}
}

func TestToFloat_NonNumericStringPassesThrough(t *testing.T) {
	got := ToFloat("abc", nil)
	if got != "abc" {
		t.Fatalf("got %v, want unchanged 'abc'", got)
	}
}
