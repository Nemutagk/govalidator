package validate

import "testing"

func TestRequiredWithout_NoOptions(t *testing.T) {
	errors := make(map[string]interface{})
	errors, skip := RequiredWithout("name", nil, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_without")
	if !found || msgs[0] != "La opción no está definida" {
		t.Fatalf("expected 'no está definida' error, got %v", errors)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWithout_TooManyOptions(t *testing.T) {
	errors := make(map[string]interface{})
	errors, skip := RequiredWithout("name", nil, map[string]any{}, []string{"a", "b"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_without")
	if !found || msgs[0] != "La opción no está definida" {
		t.Fatalf("expected 'no está definida' error, got %v", errors)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWithout_InvalidOptionsCountWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors, skip := RequiredWithout("name", nil, map[string]any{}, []string{}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_without")
	want := "La opción en la posición 2 no está definida"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWithout_ReferenceAbsentAndInputAbsentFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{}
	errors, skip := RequiredWithout("name", nil, payload, []string{"other"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_without")
	want := "El campo \"name\" debe estar definido cuando el campo \"other\" no está definido o esta vacio"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWithout_ReferenceEmptyAndInputAbsentFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"other": ""}
	errors, skip := RequiredWithout("name", nil, payload, []string{"other"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_without")
	want := "El campo \"name\" debe estar definido cuando el campo \"other\" no está definido o esta vacio"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWithout_ReferenceAbsentAndInputPresentPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"name": "Jose"}
	errors, skip := RequiredWithout("name", "Jose", payload, []string{"other"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
	if skip {
		t.Fatalf("expected skip = false, got %v", skip)
	}
}

func TestRequiredWithout_ReferencePresentPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"other": "x"}
	errors, skip := RequiredWithout("name", nil, payload, []string{"other"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
	if skip {
		t.Fatalf("expected skip = false, got %v", skip)
	}
}

func TestRequiredWithout_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{}
	errors, skip := RequiredWithout("name", nil, payload, []string{"other"}, "6", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_without")
	want := "El campo \"name\" en la posición 6 debe estar definido cuando el campo \"other\" no está definido o esta vacio"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWithout_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{}
	customErrors := map[string]string{"name.required_without": "mensaje personalizado"}
	errors, skip := RequiredWithout("name", nil, payload, []string{"other"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "name", "required_without")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}
