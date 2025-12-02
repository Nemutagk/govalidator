package validate

import "fmt"

func Nullable(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if _, ok := payload[input]; !ok {
		tmpError := fmt.Sprintf("El campo \"%s\" no existe", input)

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo \"%s\" en la posición %s no existe", input, sliceIndex)
		}

		customeErrorKey := fmt.Sprintf("%s.null", input)
		if customeError, exists := customeErrors[customeErrorKey]; exists {
			tmpError = customeError
		}

		errors = addError(input, "null", errors, tmpError)
	}

	return errors
}
