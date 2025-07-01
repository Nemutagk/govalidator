package validate

func In(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
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

	if !encontrado {
		errors = addError(input, "in", errors, "No se encontró el valor en las opciones permitidas")
	}

	return errors
}
