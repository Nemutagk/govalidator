package validate

import "testing"

func TestRequiredWithAll_NoOptions(t *testing.T) {
	errors := make(map[string]interface{})
	errors, skip := RequiredWithAll("name", nil, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_with_all")
	if !found || msgs[0] != "La opción no está definida" {
		t.Fatalf("expected 'no está definida' error, got %v", errors)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWithAll_MultipleOptionsAllDefinedAndInputMissingFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"a": "x", "b": "y"}
	errors, skip := RequiredWithAll("name", nil, payload, []string{"a", "b"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_with_all")
	want := "El campo \"name\" debe estar definido cuando los campos \"a, b\" están definidos"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWithAll_MultipleOptionsNotAllDefinedPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"a": "x"}
	errors, skip := RequiredWithAll("name", nil, payload, []string{"a", "b"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
	if skip {
		t.Fatalf("expected skip = false, got %v", skip)
	}
}

func TestRequiredWithAll_InvalidOptionsCountWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors, skip := RequiredWithAll("name", nil, map[string]any{}, []string{}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_with_all")
	want := "La opción en la posición 2 no esta definida"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWithAll_AllDefinedAndInputMissingFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"other": "x"}
	errors, skip := RequiredWithAll("name", nil, payload, []string{"other"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_with_all")
	want := "El campo \"name\" debe estar definido cuando los campos \"other\" están definidos"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWithAll_AllDefinedAndInputPresentPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"other": "x", "name": "Jose"}
	errors, skip := RequiredWithAll("name", "Jose", payload, []string{"other"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
	if skip {
		t.Fatalf("expected skip = false, got %v", skip)
	}
}

func TestRequiredWithAll_NotAllDefinedPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{}
	errors, skip := RequiredWithAll("name", nil, payload, []string{"other"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
	if skip {
		t.Fatalf("expected skip = false, got %v", skip)
	}
}

func TestRequiredWithAll_ReferenceEmptyCountsAsNotDefinedPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"other": ""}
	errors, skip := RequiredWithAll("name", nil, payload, []string{"other"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
	if skip {
		t.Fatalf("expected skip = false, got %v", skip)
	}
}

func TestRequiredWithAll_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"other": "x"}
	errors, skip := RequiredWithAll("name", nil, payload, []string{"other"}, "5", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_with_all")
	want := "El campo \"name\" en la posición 5 debe estar definido cuando los campos \"other\" están definidos"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWithAll_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"other": "x"}
	customErrors := map[string]string{"name.required_with_all": "mensaje personalizado"}
	errors, skip := RequiredWithAll("name", nil, payload, []string{"other"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "name", "required_with_all")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}
