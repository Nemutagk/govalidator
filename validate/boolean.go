package validate

import "fmt"

func Boolean(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customErrors map[string]string) map[string]interface{} {

	if _, ok := value.(bool); !ok {
		tmpError := fmt.Sprintf("El input %s no es un booleano válido", input)

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El input %s en la posición %s no es un booleano válido", input, sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.boolean", input)
		if customError, exists := customErrors[tmpErrorKey]; exists {
			tmpError = customError
		}
		errors = addError(input, "boolean", errors, tmpError)
	}

	return errors
}
