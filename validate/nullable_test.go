package validate

import "testing"

func TestNullable_PassesWhenFieldExists(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"name": "Jose"}
	errors = Nullable("name", "Jose", payload, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestNullable_PassesWhenFieldExistsWithNilValue(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"name": nil}
	errors = Nullable("name", nil, payload, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestNullable_FailsWhenFieldMissing(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Nullable("name", nil, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "null")
	want := "El campo \"name\" no existe"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestNullable_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Nullable("name", nil, map[string]any{}, []string{}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "null")
	want := "El campo \"name\" en la posición 2 no existe"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestNullable_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"name.null": "mensaje personalizado"}
	errors = Nullable("name", nil, map[string]any{}, []string{}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "name", "null")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}
