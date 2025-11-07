package validate

import "fmt"

func Confirmation(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if _, exists_input := payload[input]; !exists_input {
		fmt.Println("validate confirmation:input password does not exists")
		return errors
	}

	if payload[input] == "" {
		errors = addError(input, "confirmation", errors, "La confirmación está vacía")
		return errors
	}

	if val, exists_confirmation := payload["password_confirmation"]; !exists_confirmation && (val == nil || val == "" || val == false) {
		tmpError := "El input password_confirmation no está definido"

		customeErrorKey := fmt.Sprintf("%s.confirmation", input)
		if customeError, exists := customeErrors[customeErrorKey]; exists {
			tmpError = customeError
		}

		errors = addError(input, "confirmation", errors, tmpError)
	}

	if payload[input] != payload["password_confirmation"] {
		tmpError := "La contraseña no es igual a la confirmación"

		customeErrorKey := fmt.Sprintf("%s.confirmation", input)
		if customeError, exists := customeErrors[customeErrorKey]; exists {
			tmpError = customeError
		}

		errors = addError(input, "confirmation", errors, tmpError)
	}

	return errors
}
