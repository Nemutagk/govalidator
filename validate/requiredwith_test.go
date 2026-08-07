package validate

import "testing"

func TestRequiredWith_NoOptions(t *testing.T) {
	errors := make(map[string]interface{})
	errors, skip := RequiredWith("name", nil, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_with")
	if !found || msgs[0] != "La opción no esta definida" {
		t.Fatalf("expected 'no esta definida' error, got %v", errors)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWith_TooManyOptions(t *testing.T) {
	errors := make(map[string]interface{})
	errors, skip := RequiredWith("name", nil, map[string]any{}, []string{"a", "b", "c"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_with")
	if !found || msgs[0] != "La opción no esta definida" {
		t.Fatalf("expected 'no esta definida' error, got %v", errors)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWith_InvalidOptionsCountWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors, skip := RequiredWith("name", nil, map[string]any{}, []string{}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_with")
	want := "La opción en la posición 2 no esta definida"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWith_ReferenceFieldMissingAndInputMissing(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{}
	errors, skip := RequiredWith("name", nil, payload, []string{"other"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
	if !skip {
		t.Fatalf("expected skip = true when reference missing and input missing, got %v", skip)
	}
}

func TestRequiredWith_ReferenceFieldMissingAndInputPresent(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"name": "Jose"}
	errors, skip := RequiredWith("name", "Jose", payload, []string{"other"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
	if skip {
		t.Fatalf("expected skip = false when reference missing and input present, got %v", skip)
	}
}

func TestRequiredWith_SingleOption_ReferencePresentInputMissingFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"other": "x"}
	errors, skip := RequiredWith("name", nil, payload, []string{"other"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_with")
	want := "El campo 'name' debe estar definido cuando el campo 'other' está definido"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWith_SingleOption_ReferencePresentInputPresentPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"other": "x", "name": "Jose"}
	errors, skip := RequiredWith("name", "Jose", payload, []string{"other"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
	if skip {
		t.Fatalf("expected skip = false, got %v", skip)
	}
}

func TestRequiredWith_SingleOption_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"other": "x"}
	errors, skip := RequiredWith("name", nil, payload, []string{"other"}, "4", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_with")
	want := "El campo 'name' en la posición 4 debe estar definido cuando el campo 'other' está definido"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWith_SingleOption_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"other": "x"}
	customErrors := map[string]string{"name.required_with": "mensaje personalizado"}
	errors, skip := RequiredWith("name", nil, payload, []string{"other"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "name", "required_with")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWith_TwoOptions_ValueMatchesAndInputMissingFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"other": "active"}
	errors, skip := RequiredWith("name", nil, payload, []string{"other", "active"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_with")
	want := "El campo 'name' debe estar definido cuando el campo 'other' está definido y el valor es 'active'"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWith_TwoOptions_ValueMatchesAndInputDoesNotMatchFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"other": "active", "name": "otro-valor"}
	errors, skip := RequiredWith("name", "otro-valor", payload, []string{"other", "active"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_with")
	want := "El campo 'name' debe estar definido cuando el campo 'other' está definido y el valor es 'active'"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if skip {
		t.Fatalf("expected skip = false when input is present and non-empty, got %v", skip)
	}
}

func TestRequiredWith_TwoOptions_ValueMatchesAndInputMatchesPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"other": "active", "name": "active"}
	errors, skip := RequiredWith("name", "active", payload, []string{"other", "active"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
	if skip {
		t.Fatalf("expected skip = false, got %v", skip)
	}
}

func TestRequiredWith_TwoOptions_ValueDoesNotMatchPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"other": "inactive"}
	errors, skip := RequiredWith("name", nil, payload, []string{"other", "active"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
	if skip {
		t.Fatalf("expected skip = false, got %v", skip)
	}
}

func TestRequiredWith_TwoOptions_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"other": "active"}
	errors, skip := RequiredWith("name", nil, payload, []string{"other", "active"}, "7", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required_with")
	want := "El campo 'name' en la posición 7 debe estar definido cuando el campo 'other' está definido y el valor es 'active'"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequiredWith_TwoOptions_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"other": "active"}
	customErrors := map[string]string{"name.required_with": "mensaje personalizado"}
	errors, skip := RequiredWith("name", nil, payload, []string{"other", "active"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "name", "required_with")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}
