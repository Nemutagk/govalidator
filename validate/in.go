package validate

import "fmt"

func In(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
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
		tmpError := "No se encontró el valor en las opciones permitidas"

		customeErrorKey := fmt.Sprintf("%s.in", input)
		if customeError, exists := customeErrors[customeErrorKey]; exists {
			tmpError = customeError
		}

		errors = addError(input, "in", errors, tmpError)
	}

	return errors
}
