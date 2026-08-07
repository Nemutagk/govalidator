package validate

import "testing"

func TestPassword_FieldNotInPayloadNoError(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{}
	errors = Password("password", "whatever", payload, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestPassword_ValidPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": true}
	errors = Password("password", "Abcde1!", payload, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestPassword_NonStringTypeFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": true}
	errors = Password("password", 12345, payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "password", "type")
	if !found || len(msgs) != 1 || msgs[0] != "La contraseña debe ser una cadena de texto" {
		t.Fatalf("msgs = %v, unexpected", errors)
	}
	if len(errors) != 1 {
		t.Fatalf("expected only the type error, got %v", errors)
	}
}

func TestPassword_NonStringTypeWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": true}
	errors = Password("password", 12345, payload, []string{}, "4", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "password", "type")
	want := "La contraseña en la posición 4 debe ser una cadena de texto"
	if !found || len(msgs) != 1 || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestPassword_NonStringTypeCustomError(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": true}
	customErrors := map[string]string{"password.password": "mensaje personalizado"}
	errors = Password("password", 12345, payload, []string{}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "password", "type")
	if !found || len(msgs) != 1 || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestPassword_TooShortFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": true}
	errors = Password("password", "Ab1!2", payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "password", "min:6")
	if !found || len(msgs) != 1 || msgs[0] != "La contraseña debe tener al menos 6 caracteres" {
		t.Fatalf("msgs = %v, unexpected", errors)
	}
	if len(errors) != 1 {
		t.Fatalf("expected only the min:6 error, got %v", errors)
	}
}

func TestPassword_TooShortWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": true}
	errors = Password("password", "Ab1!2", payload, []string{}, "1", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "password", "min:6")
	want := "La contraseña en la posición 1 debe tener al menos 6 caracteres"
	if !found || len(msgs) != 1 || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestPassword_TooShortCustomError(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": true}
	customErrors := map[string]string{"password.password": "mensaje personalizado"}
	errors = Password("password", "Ab1!2", payload, []string{}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "password", "min:6")
	if !found || len(msgs) != 1 || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestPassword_MissingDigitFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": true}
	errors = Password("password", "Abcdef!", payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "password", "regex")
	if !found || len(msgs) != 1 || msgs[0] != "La contraseña debe contener al menos un número" {
		t.Fatalf("msgs = %v, unexpected", errors)
	}
}

func TestPassword_MissingDigitWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": true}
	errors = Password("password", "Abcdef!", payload, []string{}, "5", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "password", "regex")
	want := "La contraseña en la posición 5 debe contener al menos un número"
	if !found || len(msgs) != 1 || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestPassword_MissingLowercaseFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": true}
	errors = Password("password", "ABCDEF1!", payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "password", "regex")
	if !found || len(msgs) != 1 || msgs[0] != "La contraseña debe contener al menos una letra minúscula" {
		t.Fatalf("msgs = %v, unexpected", errors)
	}
}

func TestPassword_MissingUppercaseFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": true}
	errors = Password("password", "abcdef1!", payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "password", "regex")
	if !found || len(msgs) != 1 || msgs[0] != "La contraseña debe contener al menos una letra mayúscula" {
		t.Fatalf("msgs = %v, unexpected", errors)
	}
}

func TestPassword_MissingSpecialCharFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": true}
	errors = Password("password", "Abcdef12", payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "password", "regex")
	want := "La contraseña debe contener al menos un carácter especial ($#%&/()!_-)"
	if !found || len(msgs) != 1 || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestPassword_MissingSpecialCharWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": true}
	errors = Password("password", "Abcdef12", payload, []string{}, "7", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "password", "regex")
	want := "La contraseña en la posición 7 debe contener al menos un carácter especial ($#%&/()!_-)"
	if !found || len(msgs) != 1 || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestPassword_RegexCustomError(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": true}
	customErrors := map[string]string{"password.password": "mensaje personalizado"}
	errors = Password("password", "Abcdef12", payload, []string{}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "password", "regex")
	if !found || len(msgs) != 1 || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestPassword_MultipleFailuresAccumulateInOrder(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": true}
	errors = Password("password", "abcdef", payload, []string{}, "", errors, testAddError, map[string]string{})

	if _, found := getErrorMsgs(errors, "password", "min:6"); found {
		t.Fatalf("did not expect min:6 error, got %v", errors)
	}

	msgs, found := getErrorMsgs(errors, "password", "regex")
	want := []string{
		"La contraseña debe contener al menos un número",
		"La contraseña debe contener al menos una letra mayúscula",
		"La contraseña debe contener al menos un carácter especial ($#%&/()!_-)",
	}
	if !found || len(msgs) != len(want) {
		t.Fatalf("msgs = %v, want %v", msgs, want)
	}
	for i := range want {
		if msgs[i] != want[i] {
			t.Fatalf("msgs[%d] = %q, want %q", i, msgs[i], want[i])
		}
	}
}

func TestPassword_UsesActualFieldNameNotHardcoded(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"new_password": true}
	errors = Password("new_password", "abc", payload, []string{}, "", errors, testAddError, map[string]string{})

	if _, found := getErrorMsgs(errors, "password", "min:6"); found {
		t.Fatalf("expected no errors under the literal 'password' key, got %v", errors)
	}

	msgs, found := getErrorMsgs(errors, "new_password", "min:6")
	if !found || len(msgs) != 1 || msgs[0] != "La contraseña debe tener al menos 6 caracteres" {
		t.Fatalf("expected error under 'new_password' key, got %v", errors)
	}
}

func TestPassword_AllRulesFail(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": true}
	errors = Password("password", "a", payload, []string{}, "", errors, testAddError, map[string]string{})

	minMsgs, found := getErrorMsgs(errors, "password", "min:6")
	if !found || len(minMsgs) != 1 || minMsgs[0] != "La contraseña debe tener al menos 6 caracteres" {
		t.Fatalf("min:6 msgs = %v, unexpected", errors)
	}

	regexMsgs, found := getErrorMsgs(errors, "password", "regex")
	want := []string{
		"La contraseña debe contener al menos un número",
		"La contraseña debe contener al menos una letra mayúscula",
		"La contraseña debe contener al menos un carácter especial ($#%&/()!_-)",
	}
	if !found || len(regexMsgs) != len(want) {
		t.Fatalf("regex msgs = %v, want %v", regexMsgs, want)
	}
	for i := range want {
		if regexMsgs[i] != want[i] {
			t.Fatalf("regexMsgs[%d] = %q, want %q", i, regexMsgs[i], want[i])
		}
	}
}
