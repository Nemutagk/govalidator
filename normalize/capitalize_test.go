package normalize

import "testing"

func TestCapitalize_UppercasesFirstRuneOnly(t *testing.T) {
	got := Capitalize("juan perez", nil)
	if got != "Juan perez" {
		t.Fatalf("got %q, want %q", got, "Juan perez")
	}
}

func TestCapitalize_EmptyStringPassesThrough(t *testing.T) {
	got := Capitalize("", nil)
	if got != "" {
		t.Fatalf("got %q, want empty string", got)
	}
}

func TestCapitalize_NonStringPassesThrough(t *testing.T) {
	got := Capitalize(42, nil)
	if got != 42 {
		t.Fatalf("got %v, want 42", got)
	}
}
