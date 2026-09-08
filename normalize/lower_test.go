package normalize

import "testing"

func TestLower_TransformsString(t *testing.T) {
	got := Lower("HeLLo", nil)
	if got != "hello" {
		t.Fatalf("got %v, want %q", got, "hello")
	}
}

func TestLower_NonStringPassesThrough(t *testing.T) {
	got := Lower(42, nil)
	if got != 42 {
		t.Fatalf("got %v, want 42", got)
	}
}

func TestLower_NilPassesThrough(t *testing.T) {
	got := Lower(nil, nil)
	if got != nil {
		t.Fatalf("got %v, want nil", got)
	}
}
