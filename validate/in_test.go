package validate

import "testing"

func TestIn_NilValue(t *testing.T) {
	errors := make(map[string]interface{})
	errors = In("status", nil, map[string]any{}, []string{"active", "inactive"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestIn_EmptyStringValue(t *testing.T) {
	errors := make(map[string]interface{})
	errors = In("status", "", map[string]any{}, []string{"active", "inactive"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestIn_ValuePasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = In("status", "active", map[string]any{}, []string{"active", "inactive"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestIn_ValueFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = In("status", "pending", map[string]any{}, []string{"active", "inactive"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "status", "in")
	if !found || msgs[0] != "No se encontró el valor en las opciones permitidas" {
		t.Fatalf("expected 'no encontrado' error, got %v", errors)
	}
}

func TestIn_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = In("status", "pending", map[string]any{}, []string{"active", "inactive"}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "status", "in")
	want := "El valor en la posición 2 no se encontró en las opciones permitidas"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestIn_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"status.in": "mensaje personalizado"}
	errors = In("status", "pending", map[string]any{}, []string{"active", "inactive"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "status", "in")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestIn_NoOptionsFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = In("status", "active", map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "status", "in")
	if !found || msgs[0] != "No se encontró el valor en las opciones permitidas" {
		t.Fatalf("expected 'no encontrado' error, got %v", errors)
	}
}
