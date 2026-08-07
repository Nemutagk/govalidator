package validate

import "testing"

func uniqueModels() map[string]func(data any, payload map[string]any, opts *[]string) (bool, string) {
	return map[string]func(data any, payload map[string]any, opts *[]string) (bool, string){
		"valid": func(data any, payload map[string]any, opts *[]string) (bool, string) { return true, "" },
		"invalidWithMsg": func(data any, payload map[string]any, opts *[]string) (bool, string) {
			return false, "el correo ya existe"
		},
		"invalidNoMsg": func(data any, payload map[string]any, opts *[]string) (bool, string) { return false, "" },
	}
}

func TestUnique_InvalidOptionsCount(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"email": "a@b.com"}
	errors = Unique("email", "a@b.com", payload, []string{}, "", errors, testAddError, uniqueModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "email", "unique")
	if !found || msgs[0] != "the options is not valid" {
		t.Fatalf("expected 'options is not valid' error, got %v", errors)
	}
}

func TestUnique_InputNotInPayload(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{}
	errors = Unique("email", "a@b.com", payload, []string{"valid"}, "", errors, testAddError, uniqueModels(), map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestUnique_NilModelsMap(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"email": "a@b.com"}
	errors = Unique("email", "a@b.com", payload, []string{"valid"}, "", errors, testAddError, nil, map[string]string{})

	msgs, found := getErrorMsgs(errors, "email", "unique")
	if !found || msgs[0] != "el valor no es válido" {
		t.Fatalf("expected 'valor no es válido' error, got %v", errors)
	}
}

func TestUnique_NilModelsMapWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"email": "a@b.com"}
	errors = Unique("email", "a@b.com", payload, []string{"valid"}, "2", errors, testAddError, nil, map[string]string{})

	msgs, found := getErrorMsgs(errors, "email", "unique")
	want := "el valor en la posición 2 no es válido"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestUnique_ModelKeyNotFound(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"email": "a@b.com"}
	errors = Unique("email", "a@b.com", payload, []string{"missing"}, "", errors, testAddError, uniqueModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "email", "unique")
	if !found || msgs[0] != "el valor no es válido" {
		t.Fatalf("expected 'valor no es válido' error, got %v", errors)
	}
}

func TestUnique_ValueNotString(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"email": 123}
	errors = Unique("email", 123, payload, []string{"valid"}, "", errors, testAddError, uniqueModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "email", "unique")
	if !found || msgs[0] != "el valor no es válido" {
		t.Fatalf("expected 'valor no es válido' error, got %v", errors)
	}
}

func TestUnique_ModelPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"email": "a@b.com"}
	errors = Unique("email", "a@b.com", payload, []string{"valid"}, "", errors, testAddError, uniqueModels(), map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestUnique_ModelFailsWithCustomErr(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"email": "a@b.com"}
	errors = Unique("email", "a@b.com", payload, []string{"invalidWithMsg"}, "", errors, testAddError, uniqueModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "email", "unique")
	if !found || msgs[0] != "el correo ya existe" {
		t.Fatalf("expected model custom error, got %v", errors)
	}
}

func TestUnique_ModelFailsDefaultMessage(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"email": "a@b.com"}
	errors = Unique("email", "a@b.com", payload, []string{"invalidNoMsg"}, "", errors, testAddError, uniqueModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "email", "unique")
	if !found || msgs[0] != "El valor 'a@b.com' ya está registrado" {
		t.Fatalf("expected default already registered error, got %v", errors)
	}
}

func TestUnique_ModelFailsCustomeErrorsOverride(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"email": "a@b.com"}
	customErrors := map[string]string{"email.unique": "mensaje personalizado"}
	errors = Unique("email", "a@b.com", payload, []string{"invalidNoMsg"}, "", errors, testAddError, uniqueModels(), customErrors)

	msgs, found := getErrorMsgs(errors, "email", "unique")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestUnique_ModelFailsCustomErrTakesPrecedenceOverCustomeErrors(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"email": "a@b.com"}
	customErrors := map[string]string{"email.unique": "mensaje personalizado"}
	errors = Unique("email", "a@b.com", payload, []string{"invalidWithMsg"}, "", errors, testAddError, uniqueModels(), customErrors)

	msgs, found := getErrorMsgs(errors, "email", "unique")
	if !found || msgs[0] != "el correo ya existe" {
		t.Fatalf("expected model error to take precedence, got %v", errors)
	}
}
