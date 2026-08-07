package validate

import "testing"

func TestRequired_MissingField(t *testing.T) {
	errors := make(map[string]interface{})
	errors, skip := Required("name", nil, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required")
	if !found || len(msgs) != 1 || msgs[0] != "El campo name no está definido" {
		t.Fatalf("expected 'no está definido' error, got %v", errors)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequired_EmptyStringValue(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"name": ""}
	errors, skip := Required("name", "", payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required")
	if !found || len(msgs) != 1 || msgs[0] != "El campo name está vacío" {
		t.Fatalf("expected 'está vacío' error, got %v", errors)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequired_NilValue(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"name": nil}
	errors, skip := Required("name", nil, payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required")
	if !found || msgs[0] != "El campo name está vacío" {
		t.Fatalf("expected 'está vacío' error, got %v", errors)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequired_ValuePresentPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"name": "Jose"}
	errors, skip := Required("name", "Jose", payload, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
	if skip {
		t.Fatalf("expected skip = false, got %v", skip)
	}
}

func TestRequired_MissingWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors, skip := Required("name", nil, map[string]any{}, []string{}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required")
	want := "El campo en la posición 2 no está definido"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequired_EmptyWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"name": ""}
	errors, skip := Required("name", "", payload, []string{}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "required")
	want := "El campo en la posición 2 no está definido"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequired_CustomErrorMessageOnMissing(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"name.required": "mensaje personalizado"}
	errors, skip := Required("name", nil, map[string]any{}, []string{}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "name", "required")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}

func TestRequired_CustomErrorMessageOnEmpty(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"name.required": "mensaje personalizado"}
	payload := map[string]any{"name": ""}
	errors, skip := Required("name", "", payload, []string{}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "name", "required")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
	if !skip {
		t.Fatalf("expected skip = true, got %v", skip)
	}
}
