package validate

func NotIn(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	value := payload[input]

	if value == nil || value == "" {
		return errors
	}

	encontrado := true
	for _, option := range options {
		if option == value {
			encontrado = false
			break
		}
	}

	if encontrado {
		errors = addError(input, "not_in", errors, "Se encontró el valor en las opciones no permitidas")
	}

	return errors
}
