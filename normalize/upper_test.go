package normalize

import "testing"

func TestUpper_TransformsString(t *testing.T) {
	got := Upper("HeLLo", nil)
	if got != "HELLO" {
		t.Fatalf("got %v, want %q", got, "HELLO")
	}
}

func TestUpper_NonStringPassesThrough(t *testing.T) {
	got := Upper(42, nil)
	if got != 42 {
		t.Fatalf("got %v, want 42", got)
	}
}
