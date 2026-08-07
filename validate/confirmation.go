package validate

import "fmt"

func Confirmation(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if _, exists_input := payload[input]; !exists_input {
		return errors
	}

	if payload[input] == "" {
		tmpError := fmt.Sprintf("El campo %s está vacío", input)

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo %s en la posición %s está vacío", input, sliceIndex)
		}

		errors = addError(input, "confirmation", errors, tmpError)
		return errors
	}

	confirmationField := input + "_confirmation"

	if val, exists_confirmation := payload[confirmationField]; !exists_confirmation && (val == nil || val == "" || val == false) {
		tmpError := fmt.Sprintf("El campo %s no está definido", confirmationField)

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo %s en la posición %s no está definido", confirmationField, sliceIndex)
		}

		customeErrorKey := fmt.Sprintf("%s.confirmation", input)
		if customeError, exists := customeErrors[customeErrorKey]; exists {
			tmpError = customeError
		}

		errors = addError(input, "confirmation", errors, tmpError)
	}

	if payload[input] != payload[confirmationField] {
		tmpError := fmt.Sprintf("El campo %s no coincide con %s", input, confirmationField)

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo %s en la posición %s no coincide con %s", input, sliceIndex, confirmationField)
		}

		customeErrorKey := fmt.Sprintf("%s.confirmation", input)
		if customeError, exists := customeErrors[customeErrorKey]; exists {
			tmpError = customeError
		}

		errors = addError(input, "confirmation", errors, tmpError)
	}

	return errors
}
