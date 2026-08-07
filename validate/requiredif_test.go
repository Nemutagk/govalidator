package validate

import "testing"

func TestRequiredIf_NoOptions(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{}
	errors = RequiredIf("name", nil, body, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_if")
	want := "La regla required_if requiere al menos 1 parámetro (ruta del nodo)"
	if !found || len(msgs) != 1 || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestRequiredIf_NodeMissingPasses(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{}
	errors = RequiredIf("name", nil, body, []string{"other"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestRequiredIf_NodePresentValuePresentPasses(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"other": "x"}
	errors = RequiredIf("name", "Jose", body, []string{"other"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestRequiredIf_NodePresentValueMissingFails(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"other": "x"}
	errors = RequiredIf("name", nil, body, []string{"other"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_if")
	want := "El campo name es requerido cuando other está presente"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestRequiredIf_NodePresentValueEmptyStringFails(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"other": "x"}
	errors = RequiredIf("name", "", body, []string{"other"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_if")
	want := "El campo name es requerido cuando other está presente"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestRequiredIf_NestedNodePath(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"parent": map[string]any{"child": "x"}}
	errors = RequiredIf("name", nil, body, []string{"parent.child"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_if")
	want := "El campo name es requerido cuando parent.child está presente"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestRequiredIf_ExpectedValueMatchesAndValueMissingFails(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"other": "active"}
	errors = RequiredIf("name", nil, body, []string{"other", "active"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_if")
	want := "El campo name es requerido cuando other es active"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestRequiredIf_ExpectedValueDoesNotMatchPasses(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"other": "inactive"}
	errors = RequiredIf("name", nil, body, []string{"other", "active"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestRequiredIf_ExpectedValueMatchesAndValuePresentPasses(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"other": "active"}
	errors = RequiredIf("name", "Jose", body, []string{"other", "active"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestRequiredIf_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"other": "x"}
	errors = RequiredIf("name", nil, body, []string{"other"}, "3", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_if")
	want := "El campo en la posición 3 es requerido cuando other está presente"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestRequiredIf_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"other": "x"}
	customErrors := map[string]string{"name.required_if": "mensaje personalizado"}
	errors = RequiredIf("name", nil, body, []string{"other"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "name", "required_if")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}
