package validate

import "testing"

func TestNotIn_NilValue(t *testing.T) {
	errors := make(map[string]interface{})
	errors = NotIn("status", nil, map[string]any{}, []string{"banned", "blocked"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestNotIn_EmptyStringValue(t *testing.T) {
	errors := make(map[string]interface{})
	errors = NotIn("status", "", map[string]any{}, []string{"banned", "blocked"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestNotIn_ValuePasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = NotIn("status", "active", map[string]any{}, []string{"banned", "blocked"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestNotIn_ValueFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = NotIn("status", "banned", map[string]any{}, []string{"banned", "blocked"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "status", "notin")
	if !found || msgs[0] != "El valor se encontró en las opciones prohibidas" {
		t.Fatalf("expected 'encontrado en prohibidas' error, got %v", errors)
	}
}

func TestNotIn_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = NotIn("status", "banned", map[string]any{}, []string{"banned", "blocked"}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "status", "notin")
	want := "El valor en la posición 2 se encontró en las opciones prohibidas"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestNotIn_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"status.notin": "mensaje personalizado"}
	errors = NotIn("status", "banned", map[string]any{}, []string{"banned", "blocked"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "status", "notin")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestNotIn_NoOptionsPasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = NotIn("status", "active", map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}
