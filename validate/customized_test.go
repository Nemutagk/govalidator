package validate

import "testing"

func customizedModels() map[string]func(data any, payload map[string]any, opts *[]string) (bool, string) {
	return map[string]func(data any, payload map[string]any, opts *[]string) (bool, string){
		"valid": func(data any, payload map[string]any, opts *[]string) (bool, string) { return true, "" },
		"invalidWithMsg": func(data any, payload map[string]any, opts *[]string) (bool, string) {
			return false, "la edad debe ser mayor a 18"
		},
		"invalidNoMsg": func(data any, payload map[string]any, opts *[]string) (bool, string) { return false, "" },
	}
}

func TestCustomized_NoOptions(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Customized("age", 15, map[string]any{}, []string{}, "", errors, testAddError, customizedModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "customized")
	if !found || msgs[0] != "La función de validación personalizada no está definida" {
		t.Fatalf("expected 'no definida' error, got %v", errors)
	}
}

func TestCustomized_NoOptionsWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Customized("age", 15, map[string]any{}, []string{}, "2", errors, testAddError, customizedModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "customized")
	want := "La función de validación personalizada en la posición 2 no está definida"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestCustomized_FunctionKeyNotFound(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Customized("age", 15, map[string]any{}, []string{"missing"}, "", errors, testAddError, customizedModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "customized")
	if !found || msgs[0] != "La función de validación personalizada no existe" {
		t.Fatalf("expected 'no existe' error, got %v", errors)
	}
}

func TestCustomized_FunctionKeyNotFoundWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Customized("age", 15, map[string]any{}, []string{"missing"}, "2", errors, testAddError, customizedModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "customized")
	want := "La función de validación personalizada en la posición 2 no existe"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestCustomized_FunctionPasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Customized("age", 20, map[string]any{}, []string{"valid"}, "", errors, testAddError, customizedModels(), map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestCustomized_FunctionFailsDefaultMessage(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Customized("age", 15, map[string]any{}, []string{"invalidNoMsg"}, "", errors, testAddError, customizedModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "customized")
	if !found || msgs[0] != "La validación personalizada ha fallado" {
		t.Fatalf("expected default failure message, got %v", errors)
	}
}

func TestCustomized_FunctionFailsDefaultMessageWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Customized("age", 15, map[string]any{}, []string{"invalidNoMsg"}, "2", errors, testAddError, customizedModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "customized")
	want := "La validación personalizada en la posición 2 ha fallado"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestCustomized_FunctionFailsWithCustomErr(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Customized("age", 15, map[string]any{}, []string{"invalidWithMsg"}, "", errors, testAddError, customizedModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "customized")
	if !found || msgs[0] != "la edad debe ser mayor a 18" {
		t.Fatalf("expected function custom error, got %v", errors)
	}
}

func TestCustomized_FunctionFailsCustomeErrorsOverride(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"age.customized": "mensaje personalizado"}
	errors = Customized("age", 15, map[string]any{}, []string{"invalidNoMsg"}, "", errors, testAddError, customizedModels(), customErrors)

	msgs, found := getErrorMsgs(errors, "age", "customized")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestCustomized_FunctionFailsCustomErrTakesPrecedenceOverCustomeErrors(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"age.customized": "mensaje personalizado"}
	errors = Customized("age", 15, map[string]any{}, []string{"invalidWithMsg"}, "", errors, testAddError, customizedModels(), customErrors)

	msgs, found := getErrorMsgs(errors, "age", "customized")
	if !found || msgs[0] != "la edad debe ser mayor a 18" {
		t.Fatalf("expected function error to take precedence, got %v", errors)
	}
}
