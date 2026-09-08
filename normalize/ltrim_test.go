package normalize

import "testing"

func TestLTrim_RemovesOnlyLeftSide(t *testing.T) {
	got := LTrim("  hello  ", nil)
	if got != "hello  " {
		t.Fatalf("got %q, want %q", got, "hello  ")
	}
}

func TestLTrim_NonStringPassesThrough(t *testing.T) {
	got := LTrim(42, nil)
	if got != 42 {
		t.Fatalf("got %v, want 42", got)
	}
}
