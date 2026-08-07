package validate

import "testing"

func TestEqual_Passes(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"status": "active"}
	errors = Equal("status", "active", payload, []string{"active"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestEqual_NoOptions(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"status": "active"}
	errors = Equal("status", "active", payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "status", "equal")
	want := "No se proporcionó un valor para comparar"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestEqual_FieldMissingFromPayload(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Equal("status", "active", map[string]any{}, []string{"active"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "status", "equal")
	want := "El campo a comparar no existe en la carga útil"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestEqual_FieldMissingFromPayloadCustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"status.equal": "mensaje personalizado"}
	errors = Equal("status", "active", map[string]any{}, []string{"active"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "status", "equal")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestEqual_FailsWhenNotEqual(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"status": "active"}
	errors = Equal("status", "inactive", payload, []string{"active"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "status", "equal")
	want := "El valor no es igual a active"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestEqual_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"status": "active"}
	errors = Equal("status", "inactive", payload, []string{"active"}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "status", "equal")
	want := "El valor en la posición 2 no es igual a active"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestEqual_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"status": "active"}
	customErrors := map[string]string{"status.equal": "mensaje personalizado"}
	errors = Equal("status", "inactive", payload, []string{"active"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "status", "equal")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}
