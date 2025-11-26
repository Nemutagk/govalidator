package validate

import (
	"fmt"
	"strings"
)

func RequiredWithoutAll(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if len(options) != 1 {
		tmpError := "La opción no está definida"

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("La opción en la posición %s no está definida", sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.required_without_all", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		return addError(input, "required_without_all", errors, tmpError)
	}

	not_defined := true
	for _, another_input := range options {
		if _, exists_another_input := payload[another_input]; !exists_another_input {
			not_defined = false
			break
		}
	}

	// fmt.Println(input+": not_defined", not_defined)

	if _, exists_input := payload[input]; !exists_input && !not_defined {
		all_inputs := strings.Join(options, ", ")
		tmpError := fmt.Sprintf("El campo \"%s\" debe estar definido cuando los campos \"%s\" no están definidos", input, all_inputs)

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo \"%s\" en la posición %s debe estar definido cuando los campos \"%s\" no están definidos", input, sliceIndex, all_inputs)
		}

		tmpErrorKey := fmt.Sprintf("%s.required_without_all", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "required_without_all", errors, tmpError)
		return errors
	}

	return errors
}
