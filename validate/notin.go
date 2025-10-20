package validate

func NotIn(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	value := payload[input]

	if value == nil || value == "" {
		return errors
	}

	encontrado := false
	for _, option := range options {
		if option == value {
			encontrado = true
			break
		}
	}

	if encontrado {
		errors = addError(input, "notin", errors, "El valor se encontró en las opciones prohibidas")
	}

	return errors
}
