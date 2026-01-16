package validate

import "fmt"

func Required(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) (map[string]interface{}, bool) {
	if _, exists_input := payload[input]; !exists_input {
		tmpError := fmt.Sprintf("El campo %s no está definido", input)

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo en la posición %s no está definido", sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.required", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "required", errors, tmpError)
		return errors, true
	}

	if payload[input] == "" || payload[input] == nil {
		tmpError := fmt.Sprintf("El campo %s está vacío", input)

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo en la posición %s no está definido", sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.required", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "required", errors, tmpError)
		return errors, true
	}

	return errors, false
}
