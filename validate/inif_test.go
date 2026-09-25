package validate

import "testing"

func TestInIf_TooFewOptionsFails(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{}
	errors = InIf("kind", "x", body, []string{"mode", "batch"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "kind", "in_if")
	want := "La regla in_if requiere al menos 3 parámetros (ruta del nodo, valor esperado y al menos un valor permitido)"
	if !found || len(msgs) != 1 || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestInIf_NodeMissingPasses(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{}
	errors = InIf("kind", "anything", body, []string{"mode", "batch", "a", "b"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestInIf_ConditionDoesNotMatchPasses(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"mode": "individual"}
	errors = InIf("kind", "not-in-any-list", body, []string{"mode", "batch", "a", "b"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestInIf_ConditionMatchesAndValueAllowedPasses(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"mode": "batch"}
	errors = InIf("kind", "a", body, []string{"mode", "batch", "a", "b"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestInIf_ConditionMatchesAndValueNotAllowedFails(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"mode": "batch"}
	errors = InIf("kind", "z", body, []string{"mode", "batch", "a", "b"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "kind", "in_if")
	want := "No se encontró el valor en las opciones permitidas"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestInIf_NilOrEmptyValuePasses(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"mode": "batch"}
	errors = InIf("kind", nil, body, []string{"mode", "batch", "a", "b"}, "", errors, testAddError, map[string]string{})
	if len(errors) != 0 {
		t.Fatalf("expected no errors for nil value, got %v", errors)
	}

	errors = InIf("kind", "", body, []string{"mode", "batch", "a", "b"}, "", errors, testAddError, map[string]string{})
	if len(errors) != 0 {
		t.Fatalf("expected no errors for empty value, got %v", errors)
	}
}

func TestInIf_NestedNodePath(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"parent": map[string]any{"mode": "batch"}}
	errors = InIf("kind", "z", body, []string{"parent.mode", "batch", "a", "b"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "kind", "in_if")
	want := "No se encontró el valor en las opciones permitidas"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestInIf_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"mode": "batch"}
	errors = InIf("kind", "z", body, []string{"mode", "batch", "a", "b"}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "kind", "in_if")
	want := "El valor en la posición 2 no se encontró en las opciones permitidas"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestInIf_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"mode": "batch"}
	customErrors := map[string]string{"kind.in_if": "mensaje personalizado"}
	errors = InIf("kind", "z", body, []string{"mode", "batch", "a", "b"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "kind", "in_if")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}
