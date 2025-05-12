package validate

func Boolean(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	if _, ok := payload[input]; !ok {
		return errors
	}

	value := payload[input]

	if _, ok := value.(bool); !ok {
		errors = addError(input, "boolean", errors, "The input "+input+" is not a boolean")
	}

	return errors
}
