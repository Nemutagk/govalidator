package validate

import "testing"

func TestBoolean_Passes(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Boolean("active", true, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestBoolean_FailsWithString(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Boolean("active", "true", map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "active", "boolean")
	want := "El input active no es un booleano válido"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestBoolean_FailsWithNil(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Boolean("active", nil, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	if _, found := getErrorMsgs(errors, "active", "boolean"); !found {
		t.Fatalf("expected nil value to fail, got %v", errors)
	}
}

func TestBoolean_FailsWithInt(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Boolean("active", 1, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	if _, found := getErrorMsgs(errors, "active", "boolean"); !found {
		t.Fatalf("expected int value to fail, got %v", errors)
	}
}

func TestBoolean_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Boolean("active", "true", map[string]any{}, []string{}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "active", "boolean")
	want := "El input active en la posición 2 no es un booleano válido"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestBoolean_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"active.boolean": "mensaje personalizado"}
	errors = Boolean("active", "true", map[string]any{}, []string{}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "active", "boolean")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}
