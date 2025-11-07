package validate

import "fmt"

func Boolean(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customErrors map[string]string) map[string]interface{} {
	if _, ok := payload[input]; !ok {
		return errors
	}

	if _, ok := payload[input].(bool); !ok {
		tmpError := fmt.Sprintf("El input %s no es un booleano válido", input)
		tmpErrorKey := fmt.Sprintf("%s.boolean", input)
		if customError, exists := customErrors[tmpErrorKey]; exists {
			tmpError = customError
		}
		errors = addError(input, "boolean", errors, tmpError)
	}

	return errors
}
