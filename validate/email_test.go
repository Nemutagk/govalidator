package validate

import "testing"

func TestEmail_ValidPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"email": "user@example.com"}
	errors = Email("email", nil, payload, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestEmail_ValidWithSubdomainPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"email": "user+tag@mail.example.co"}
	errors = Email("email", nil, payload, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestEmail_InvalidFormatFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"email": "not-an-email"}
	errors = Email("email", nil, payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "email", "email")
	if !found || len(msgs) != 1 || msgs[0] != "El campo no es un correo electrónico válido" {
		t.Fatalf("expected invalid email error, got %v", errors)
	}
}

func TestEmail_ShortTldFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"email": "user@example.c"}
	errors = Email("email", nil, payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "email", "email")
	if !found || len(msgs) != 1 || msgs[0] != "El campo no es un correo electrónico válido" {
		t.Fatalf("expected invalid email error, got %v", errors)
	}
}

func TestEmail_InvalidWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"email": "not-an-email"}
	errors = Email("email", nil, payload, []string{}, "3", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "email", "email")
	want := "El campo en la posición 3 no es un correo electrónico válido"
	if !found || len(msgs) != 1 || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestEmail_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"email": "not-an-email"}
	customErrors := map[string]string{"email.email": "mensaje personalizado"}
	errors = Email("email", nil, payload, []string{}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "email", "email")
	if !found || len(msgs) != 1 || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestEmail_NonStringValueNoError(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"email": 123}
	errors = Email("email", nil, payload, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors for non-string payload value, got %v", errors)
	}
}

func TestEmail_MissingFieldNoError(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{}
	errors = Email("email", nil, payload, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors for missing field, got %v", errors)
	}
}
