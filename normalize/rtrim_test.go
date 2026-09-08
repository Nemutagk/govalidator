package normalize

import "testing"

func TestRTrim_RemovesOnlyRightSide(t *testing.T) {
	got := RTrim("  hello  ", nil)
	if got != "  hello" {
		t.Fatalf("got %q, want %q", got, "  hello")
	}
}

func TestRTrim_NonStringPassesThrough(t *testing.T) {
	got := RTrim(42, nil)
	if got != 42 {
		t.Fatalf("got %v, want 42", got)
	}
}
