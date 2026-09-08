package normalize

import "testing"

func TestTrimChar_RemovesGivenCutset(t *testing.T) {
	got := TrimChar("---hello---", []string{"-"})
	if got != "hello" {
		t.Fatalf("got %q, want %q", got, "hello")
	}
}

func TestTrimChar_NoOptionsPassesThrough(t *testing.T) {
	got := TrimChar("---hello---", nil)
	if got != "---hello---" {
		t.Fatalf("got %q, want unchanged value", got)
	}
}

func TestTrimChar_NonStringPassesThrough(t *testing.T) {
	got := TrimChar(42, []string{"-"})
	if got != 42 {
		t.Fatalf("got %v, want 42", got)
	}
}
