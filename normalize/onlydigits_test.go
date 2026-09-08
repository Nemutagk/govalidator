package normalize

import "testing"

func TestOnlyDigits_StripsNonDigits(t *testing.T) {
	got := OnlyDigits("+52 (55) 1234-5678", nil)
	if got != "525512345678" {
		t.Fatalf("got %q, want %q", got, "525512345678")
	}
}

func TestOnlyDigits_NonStringPassesThrough(t *testing.T) {
	got := OnlyDigits(42, nil)
	if got != 42 {
		t.Fatalf("got %v, want 42", got)
	}
}
