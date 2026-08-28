package validate

import "testing"

func TestRequiredIfAll_NoOptions(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{}
	errors, _ = RequiredIfAll("name", nil, body, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_if_all")
	want := "La regla required_if_all requiere un número par de parámetros (pares de ruta del nodo y valor esperado)"
	if !found || len(msgs) != 1 || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestRequiredIfAll_SingleOption(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{}
	errors, _ = RequiredIfAll("name", nil, body, []string{"other"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_if_all")
	want := "La regla required_if_all requiere un número par de parámetros (pares de ruta del nodo y valor esperado)"
	if !found || len(msgs) != 1 || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestRequiredIfAll_OddOptions(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"other": "active", "other2": "active2"}
	errors, _ = RequiredIfAll("name", nil, body, []string{"other", "active", "other2"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_if_all")
	want := "La regla required_if_all requiere un número par de parámetros (pares de ruta del nodo y valor esperado)"
	if !found || len(msgs) != 1 || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestRequiredIfAll_OnePairMatchesAndValueMissingFails(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"other": "active"}
	errors, _ = RequiredIfAll("name", nil, body, []string{"other", "active"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_if_all")
	want := "El campo name es requerido cuando other es active"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestRequiredIfAll_OnePairMatchesAndValuePresentPasses(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"other": "active"}
	errors, _ = RequiredIfAll("name", "Jose", body, []string{"other", "active"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestRequiredIfAll_AllPairsMatchAndValueMissingFails(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"other": "active", "other2": "active2"}
	errors, _ = RequiredIfAll("name", nil, body, []string{"other", "active", "other2", "active2"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_if_all")
	want := "El campo name es requerido cuando other es active y other2 es active2"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestRequiredIfAll_AllPairsMatchAndValueEmptyStringFails(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"other": "active", "other2": "active2"}
	errors, _ = RequiredIfAll("name", "", body, []string{"other", "active", "other2", "active2"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_if_all")
	want := "El campo name es requerido cuando other es active y other2 es active2"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestRequiredIfAll_OnePairFailsToMatchPasses(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"other": "active", "other2": "inactive"}
	errors, _ = RequiredIfAll("name", nil, body, []string{"other", "active", "other2", "active2"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestRequiredIfAll_OneNodeMissingPasses(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"other": "active"}
	errors, _ = RequiredIfAll("name", nil, body, []string{"other", "active", "other2", "active2"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestRequiredIfAll_AllPairsMatchAndValuePresentPasses(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"other": "active", "other2": "active2"}
	errors, _ = RequiredIfAll("name", "Jose", body, []string{"other", "active", "other2", "active2"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestRequiredIfAll_NestedNodePath(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"parent": map[string]any{"child": "x"}}
	errors, _ = RequiredIfAll("name", nil, body, []string{"parent.child", "x"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_if_all")
	want := "El campo name es requerido cuando parent.child es x"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestRequiredIfAll_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"other": "active", "other2": "active2"}
	errors, _ = RequiredIfAll("name", nil, body, []string{"other", "active", "other2", "active2"}, "3", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_if_all")
	want := "El campo en la posición 3 es requerido cuando other es active y other2 es active2"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestRequiredIfAll_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	body := map[string]any{"other": "active", "other2": "active2"}
	customErrors := map[string]string{"name.required_if_all": "mensaje personalizado"}
	errors, _ = RequiredIfAll("name", nil, body, []string{"other", "active", "other2", "active2"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "name", "required_if_all")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}
