package validate

func Nullable(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) (map[string]interface{}, bool) {
	if _, exists := payload[input]; !exists {
		addError(input, "nullable", errors, "El valor no existe.")
		return errors, false
	}

	if payload[input] != nil {
		return errors, false
	}

	return errors, true
}
