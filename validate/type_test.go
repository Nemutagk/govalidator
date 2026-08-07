package validate

import "testing"

func TestType_SkipsWhenFieldMissingFromPayload(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Type("age", nil, map[string]any{}, []string{"int"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestType_Passes(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"age": 5}
	errors = Type("age", 5, payload, []string{"int"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestType_NoOptions(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"age": 5}
	errors = Type("age", 5, payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "type")
	want := "El tipo no está definido"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestType_NoOptionsWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"age": 5}
	errors = Type("age", 5, payload, []string{}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "type")
	want := "El tipo en la posición 2 no está definido"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestType_NoOptionsCustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"age": 5}
	customErrors := map[string]string{"age.type": "mensaje personalizado"}
	errors = Type("age", 5, payload, []string{}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "age", "type")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestType_NilValueWithNullableOptionPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"age": nil}
	errors = Type("age", nil, payload, []string{"int", "nullable"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestType_NilValueWithoutNullableFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"age": nil}
	errors = Type("age", nil, payload, []string{"int"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "type")
	want := "El campo \"age\" no puede ser nulo"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestType_NilValueWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"age": nil}
	errors = Type("age", nil, payload, []string{"int"}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "type")
	want := "El campo \"age\" en la posición 2 no puede ser nulo"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestType_NilValueCustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"age": nil}
	customErrors := map[string]string{"age.type": "mensaje personalizado"}
	errors = Type("age", nil, payload, []string{"int"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "age", "type")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestType_MismatchFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"age": "5"}
	errors = Type("age", "5", payload, []string{"int"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "type")
	want := "El tipo del campo \"age\" no es \"int\""
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestType_MismatchWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"age": "5"}
	errors = Type("age", "5", payload, []string{"int"}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "type")
	want := "El tipo en la posición 2 del campo \"age\" no es \"int\""
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestType_MismatchCustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"age": "5"}
	customErrors := map[string]string{"age.type": "mensaje personalizado"}
	errors = Type("age", "5", payload, []string{"int"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "age", "type")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}
