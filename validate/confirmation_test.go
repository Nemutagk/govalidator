package validate

import "testing"

func TestConfirmation_FieldNotInPayloadNoError(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{}
	errors = Confirmation("password", nil, payload, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestConfirmation_EmptyValueFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": ""}
	errors = Confirmation("password", "", payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "password", "confirmation")
	if !found || len(msgs) != 1 || msgs[0] != "El campo password está vacío" {
		t.Fatalf("msgs = %v, unexpected", errors)
	}
}

func TestConfirmation_EmptyValueWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": ""}
	errors = Confirmation("password", "", payload, []string{}, "4", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "password", "confirmation")
	want := "El campo password en la posición 4 está vacío"
	if !found || len(msgs) != 1 || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestConfirmation_MissingConfirmationFieldFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": "secret123"}
	errors = Confirmation("password", "secret123", payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "password", "confirmation")
	want := []string{
		"El campo password_confirmation no está definido",
		"El campo password no coincide con password_confirmation",
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

func TestConfirmation_MissingConfirmationFieldWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": "secret123"}
	errors = Confirmation("password", "secret123", payload, []string{}, "3", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "password", "confirmation")
	want := []string{
		"El campo password_confirmation en la posición 3 no está definido",
		"El campo password en la posición 3 no coincide con password_confirmation",
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

func TestConfirmation_MismatchFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": "secret123", "password_confirmation": "other"}
	errors = Confirmation("password", "secret123", payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "password", "confirmation")
	if !found || len(msgs) != 1 || msgs[0] != "El campo password no coincide con password_confirmation" {
		t.Fatalf("msgs = %v, unexpected", errors)
	}
}

func TestConfirmation_MismatchWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": "secret123", "password_confirmation": "other"}
	errors = Confirmation("password", "secret123", payload, []string{}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "password", "confirmation")
	want := "El campo password en la posición 2 no coincide con password_confirmation"
	if !found || len(msgs) != 1 || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestConfirmation_MatchPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": "secret123", "password_confirmation": "secret123"}
	errors = Confirmation("password", "secret123", payload, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestConfirmation_MismatchCustomError(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": "secret123", "password_confirmation": "other"}
	customErrors := map[string]string{"password.confirmation": "mensaje personalizado"}
	errors = Confirmation("password", "secret123", payload, []string{}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "password", "confirmation")
	if !found || len(msgs) != 1 || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestConfirmation_MissingConfirmationFieldCustomError(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"password": "secret123"}
	customErrors := map[string]string{"password.confirmation": "mensaje personalizado"}
	errors = Confirmation("password", "secret123", payload, []string{}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "password", "confirmation")
	want := []string{"mensaje personalizado", "mensaje personalizado"}
	if !found || len(msgs) != len(want) {
		t.Fatalf("msgs = %v, want %v", msgs, want)
	}
	for i := range want {
		if msgs[i] != want[i] {
			t.Fatalf("msgs[%d] = %q, want %q", i, msgs[i], want[i])
		}
	}
}

func TestConfirmation_UsesFieldSpecificConfirmationFieldPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"pin": "1234", "pin_confirmation": "1234"}
	errors = Confirmation("pin", "1234", payload, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestConfirmation_UsesFieldSpecificConfirmationFieldMismatchFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"pin": "1234", "pin_confirmation": "5678"}
	errors = Confirmation("pin", "1234", payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "pin", "confirmation")
	if !found || len(msgs) != 1 || msgs[0] != "El campo pin no coincide con pin_confirmation" {
		t.Fatalf("msgs = %v, unexpected", errors)
	}
}

func TestConfirmation_UsesFieldSpecificConfirmationFieldMissingFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"pin": "1234"}
	errors = Confirmation("pin", "1234", payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "pin", "confirmation")
	want := []string{
		"El campo pin_confirmation no está definido",
		"El campo pin no coincide con pin_confirmation",
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
