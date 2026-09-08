package normalize

import "testing"

func TestRemoveSpaces_RemovesAllWhitespace(t *testing.T) {
	got := RemoveSpaces("12 34  56\t78", nil)
	if got != "12345678" {
		t.Fatalf("got %q, want %q", got, "12345678")
	}
}

func TestRemoveSpaces_NonStringPassesThrough(t *testing.T) {
	got := RemoveSpaces(42, nil)
	if got != 42 {
		t.Fatalf("got %v, want 42", got)
	}
}
