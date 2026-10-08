package validate

import (
	"fmt"
	"testing"
)

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
	if !found || msgs[0] != "la configuración de conexión no es válida" {
		t.Fatalf("expected 'configuración no es válida' error, got %v", errors)
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

func TestUnique_ValueNotScalar(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"email": []string{"x"}}
	errors = Unique("email", []string{"x"}, payload, []string{"valid"}, "", errors, testAddError, uniqueModels(), map[string]string{})

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

func TestUnique_NumericValueReachesModelAsString(t *testing.T) {
	var received any
	models := map[string]func(data any, payload map[string]any, opts *[]string) (bool, string){
		"capture": func(data any, payload map[string]any, opts *[]string) (bool, string) {
			received = data
			return true, ""
		},
	}

	for _, v := range []any{int64(123), 7, uint8(3), 1.5, true} {
		want := fmt.Sprint(v)
		errors := make(map[string]interface{})
		errors = Unique("id", v, map[string]any{"id": v}, []string{"capture"}, "", errors, testAddError, models, map[string]string{})

		if len(errors) != 0 {
			t.Fatalf("value %v (%T): expected no errors, got %v", v, v, errors)
		}
		if received != want {
			t.Fatalf("model received %v (%T), want the string %q", received, received, want)
		}
	}
}

func TestUnique_NumericValueAlreadyRegisteredMessage(t *testing.T) {
	models := map[string]func(data any, payload map[string]any, opts *[]string) (bool, string){
		"taken": func(data any, payload map[string]any, opts *[]string) (bool, string) { return false, "" },
	}

	errors := make(map[string]interface{})
	errors = Unique("id", int64(123), map[string]any{"id": int64(123)}, []string{"taken"}, "", errors, testAddError, models, map[string]string{})

	msgs, found := getErrorMsgs(errors, "id", "unique")
	want := "El valor '123' ya está registrado"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}
