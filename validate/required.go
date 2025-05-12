package validate

func Required(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	if _, exists_input := payload[input]; !exists_input {
		errors = addError(input, "required", errors, "The input "+input+" is not defined")
		return errors
	}

	if payload[input] == "" || payload[input] == nil {
		errors = addError(input, "required", errors, "The input "+input+" is empty")
		return errors
	}

	return errors
}
