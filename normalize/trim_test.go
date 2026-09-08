package normalize

import "testing"

func TestTrim_RemovesBothSides(t *testing.T) {
	got := Trim("  hello  ", nil)
	if got != "hello" {
		t.Fatalf("got %q, want %q", got, "hello")
	}
}

func TestTrim_NonStringPassesThrough(t *testing.T) {
	got := Trim([]int{1, 2}, nil)
	if got.([]int)[0] != 1 {
		t.Fatalf("expected passthrough, got %v", got)
	}
}
