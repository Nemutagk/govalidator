package validate

import "testing"

func existsModels() map[string]func(data any, payload map[string]any, opts *[]string) (bool, string) {
	return map[string]func(data any, payload map[string]any, opts *[]string) (bool, string){
		"valid": func(data any, payload map[string]any, opts *[]string) (bool, string) { return true, "" },
		"invalidWithMsg": func(data any, payload map[string]any, opts *[]string) (bool, string) {
			return false, "el usuario no está activo"
		},
		"invalidNoMsg": func(data any, payload map[string]any, opts *[]string) (bool, string) { return false, "" },
	}
}

func TestExists_InvalidOptionsCount(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Exists("user_id", "1", map[string]any{}, []string{}, "", errors, testAddError, existsModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "user_id", "exists")
	if !found || msgs[0] != "la configuración de conexión no es válida" {
		t.Fatalf("expected 'configuración no válida' error, got %v", errors)
	}
}

func TestExists_InvalidOptionsCountWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Exists("user_id", "1", map[string]any{}, []string{}, "2", errors, testAddError, existsModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "user_id", "exists")
	want := "La configuración de conexión en la posición 2 no es válida"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestExists_InvalidOptionsCountCustomError(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"user_id.exists": "mensaje personalizado"}
	errors = Exists("user_id", "1", map[string]any{}, []string{}, "", errors, testAddError, existsModels(), customErrors)

	msgs, found := getErrorMsgs(errors, "user_id", "exists")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestExists_ModelKeyNotFound(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Exists("user_id", "1", map[string]any{}, []string{"missing"}, "", errors, testAddError, existsModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "user_id", "exists")
	if !found || msgs[0] != "el modelo no está definido" {
		t.Fatalf("expected 'modelo no definido' error, got %v", errors)
	}
}

func TestExists_ModelKeyNotFoundWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Exists("user_id", "1", map[string]any{}, []string{"missing"}, "2", errors, testAddError, existsModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "user_id", "exists")
	want := "El modelo en la posición 2 no está definido"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestExists_ValueNotString(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Exists("user_id", 123, map[string]any{}, []string{"valid"}, "", errors, testAddError, existsModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "user_id", "exists")
	if !found || msgs[0] != "el valor no es válido" {
		t.Fatalf("expected 'valor no es válido' error, got %v", errors)
	}
}

func TestExists_ValueNotStringWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Exists("user_id", 123, map[string]any{}, []string{"valid"}, "2", errors, testAddError, existsModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "user_id", "exists")
	want := "El valor en la posición 2 no es válido"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestExists_ModelPasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Exists("user_id", "1", map[string]any{}, []string{"valid"}, "", errors, testAddError, existsModels(), map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestExists_ModelFailsDefaultMessage(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Exists("user_id", "1", map[string]any{}, []string{"invalidNoMsg"}, "", errors, testAddError, existsModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "user_id", "exists")
	if !found || msgs[0] != "El valor '1' no existe" {
		t.Fatalf("expected default not exists error, got %v", errors)
	}
}

func TestExists_ModelFailsDefaultMessageWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Exists("user_id", "1", map[string]any{}, []string{"invalidNoMsg"}, "2", errors, testAddError, existsModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "user_id", "exists")
	want := "El valor '1' en la posición 2 no existe"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestExists_ModelFailsWithCustomErr(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Exists("user_id", "1", map[string]any{}, []string{"invalidWithMsg"}, "", errors, testAddError, existsModels(), map[string]string{})

	msgs, found := getErrorMsgs(errors, "user_id", "exists")
	if !found || msgs[0] != "el usuario no está activo" {
		t.Fatalf("expected model custom error, got %v", errors)
	}
}

func TestExists_ModelFailsCustomeErrorsOverride(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"user_id.exists": "mensaje personalizado"}
	errors = Exists("user_id", "1", map[string]any{}, []string{"invalidNoMsg"}, "", errors, testAddError, existsModels(), customErrors)

	msgs, found := getErrorMsgs(errors, "user_id", "exists")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestExists_ModelFailsCustomErrTakesPrecedenceOverCustomeErrors(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"user_id.exists": "mensaje personalizado"}
	errors = Exists("user_id", "1", map[string]any{}, []string{"invalidWithMsg"}, "", errors, testAddError, existsModels(), customErrors)

	msgs, found := getErrorMsgs(errors, "user_id", "exists")
	if !found || msgs[0] != "el usuario no está activo" {
		t.Fatalf("expected model error to take precedence, got %v", errors)
	}
}
