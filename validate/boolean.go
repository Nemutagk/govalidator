package validate

func Boolean(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	if _, ok := payload[input]; !ok {
		return errors
	}

	if _, ok := payload[input].(bool); !ok {
		errors = addError(input, "boolean", errors, "El input "+input+" no es un booleano válido")
	}

	return errors
}
