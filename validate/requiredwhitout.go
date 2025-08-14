package validate

func RequiredWithout(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	if len(options) != 1 {
		return addError(input, "required_without", errors, "La opción no está definida")
	}

	if _, exists := payload[options[0]]; !exists {
		if _, exists_input := payload[input]; !exists_input {
			errors = addError(input, "required_without", errors, "El campo \""+input+"\" debe estar definido cuando el campo \""+options[0]+"\" no está definido")
			return errors
		}
	}

	return errors
}
