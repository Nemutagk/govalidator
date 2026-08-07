package validate

import "testing"

func TestNotEqual_Passes(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"status": "active"}
	errors = NotEqual("status", "inactive", payload, []string{"active"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestNotEqual_NoOptions(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"status": "active"}
	errors = NotEqual("status", "active", payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "status", "not_equal")
	want := "No se proporcionó un valor para comparar"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestNotEqual_FieldMissingFromPayload(t *testing.T) {
	errors := make(map[string]interface{})
	errors = NotEqual("status", "active", map[string]any{}, []string{"active"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "status", "not_equal")
	want := "El campo a comparar no existe en la carga útil"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestNotEqual_FieldMissingFromPayloadCustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"status.not_equal": "mensaje personalizado"}
	errors = NotEqual("status", "active", map[string]any{}, []string{"active"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "status", "not_equal")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestNotEqual_FailsWhenEqual(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"status": "active"}
	errors = NotEqual("status", "active", payload, []string{"active"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "status", "not_equal")
	want := "El valor es igual a active"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestNotEqual_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"status": "active"}
	errors = NotEqual("status", "active", payload, []string{"active"}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "status", "not_equal")
	want := "El valor en la posición 2 es igual a active"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestNotEqual_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"status": "active"}
	customErrors := map[string]string{"status.not_equal": "mensaje personalizado"}
	errors = NotEqual("status", "active", payload, []string{"active"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "status", "not_equal")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}
