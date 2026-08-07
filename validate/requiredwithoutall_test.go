package validate

import "testing"

func TestRequiredWithoutAll_NoOptions(t *testing.T) {
	errors := make(map[string]interface{})
	errors, skip := RequiredWithoutAll("name", nil, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_without_all")
	if !found || msgs[0] != "La opción no está definida" {
		t.Fatalf("expected 'no está definida' error, got %v", errors)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWithoutAll_MultipleOptionsAllAbsentFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{}
	errors, skip := RequiredWithoutAll("name", nil, payload, []string{"a", "b"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_without_all")
	want := "El campo \"name\" debe estar definido cuando los campos \"a, b\" no están definidos o están vacíos"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWithoutAll_MultipleOptionsOnePresentPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"a": "x"}
	errors, skip := RequiredWithoutAll("name", nil, payload, []string{"a", "b"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
	if skip {
		t.Fatalf("expected skip = false, got %v", skip)
	}
}

func TestRequiredWithoutAll_InvalidOptionsCountWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors, skip := RequiredWithoutAll("name", nil, map[string]any{}, []string{}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_without_all")
	want := "La opción en la posición 2 no está definida"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWithoutAll_ReferenceAbsentAndInputAbsentFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{}
	errors, skip := RequiredWithoutAll("name", nil, payload, []string{"other"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_without_all")
	want := "El campo \"name\" debe estar definido cuando los campos \"other\" no están definidos o están vacíos"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWithoutAll_ReferenceEmptyAndInputAbsentFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"other": ""}
	errors, skip := RequiredWithoutAll("name", nil, payload, []string{"other"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_without_all")
	want := "El campo \"name\" debe estar definido cuando los campos \"other\" no están definidos o están vacíos"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWithoutAll_ReferenceAbsentAndInputPresentPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"name": "Jose"}
	errors, skip := RequiredWithoutAll("name", "Jose", payload, []string{"other"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
	if skip {
		t.Fatalf("expected skip = false, got %v", skip)
	}
}

func TestRequiredWithoutAll_ReferencePresentPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"other": "x"}
	errors, skip := RequiredWithoutAll("name", nil, payload, []string{"other"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
	if skip {
		t.Fatalf("expected skip = false, got %v", skip)
	}
}

func TestRequiredWithoutAll_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{}
	errors, skip := RequiredWithoutAll("name", nil, payload, []string{"other"}, "8", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_without_all")
	want := "El campo \"name\" en la posición 8 debe estar definido cuando los campos \"other\" no están definidos o están vacíos"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWithoutAll_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{}
	customErrors := map[string]string{"name.required_without_all": "mensaje personalizado"}
	errors, skip := RequiredWithoutAll("name", nil, payload, []string{"other"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "name", "required_without_all")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}
