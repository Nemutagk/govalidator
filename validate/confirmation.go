package validate

import "log"

func Confirmation(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	if _, exists_input := payload[input]; !exists_input {
		log.Println("validate confirmation:input password does not exists")
		return errors
	}

	if payload[input] == "" {
		errors = addError(input, "confirmation", errors, "El campo "+input+" está vacío")
		return errors
	}

	if val, exists_confirmation := payload["password_confirmation"]; !exists_confirmation && (val == nil || val == "" || val == false) {
		errors = addError("password_confirmation", "confirmation", errors, "El campo password_confirmation no está definido")
	}

	if payload[input] != payload["password_confirmation"] {
		errors = addError(input, "confirmation", errors, "La contraseña no coincide con la confirmación")
	}

	return errors
}
