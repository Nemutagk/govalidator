package validate

import "fmt"

func NotIn(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
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
		tmpError := "El valor se encontró en las opciones prohibidas"

		customeErrorKey := fmt.Sprintf("%s.notin", input)
		if customeError, exists := customeErrors[customeErrorKey]; exists {
			tmpError = customeError
		}

		errors = addError(input, "notin", errors, tmpError)
	}

	return errors
}
