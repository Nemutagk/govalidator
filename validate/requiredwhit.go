package validate

func RequiredWith(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	if len(options) != 1 {
		return addError(input, "required_with", errors, "The options is not defined")
	}

	if _, exists_input := payload[input]; !exists_input {
		errors = addError(input, "required_with", errors, "El campo \""+input+"\" debe estar definido cuando el campo \""+options[0]+"\" está definido")
		return errors
	}

	return errors
}
