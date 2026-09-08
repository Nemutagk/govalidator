package normalize

import "testing"

func TestToBool_ParsesTruthyValues(t *testing.T) {
	for _, s := range []string{"true", "1", "yes", "on", "TRUE"} {
		if got := ToBool(s, nil); got != true {
			t.Fatalf("input %q: got %v, want true", s, got)
		}
	}
}

func TestToBool_ParsesFalsyValues(t *testing.T) {
	for _, s := range []string{"false", "0", "no", "off", "FALSE"} {
		if got := ToBool(s, nil); got != false {
			t.Fatalf("input %q: got %v, want false", s, got)
		}
	}
}

func TestToBool_UnrecognizedStringPassesThrough(t *testing.T) {
	got := ToBool("maybe", nil)
	if got != "maybe" {
		t.Fatalf("got %v, want unchanged 'maybe'", got)
	}
}

func TestToBool_NonStringPassesThrough(t *testing.T) {
	got := ToBool(1, nil)
	if got != 1 {
		t.Fatalf("got %v, want 1", got)
	}
}
