package validate

import (
	"fmt"
	"strings"

	"github.com/Nemutagk/govalidator/v2/helper"
)

func RequiredWithoutAll(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) (map[string]interface{}, bool) {
	if len(options) < 1 {
		tmpError := "La opción no está definida"

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("La opción en la posición %s no está definida", sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.required_without_all", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		return addError(input, "required_without_all", errors, tmpError), true
	}

	all_absent := true
	for _, another_input := range options {
		if validInput, exists_another_input := payload[another_input]; exists_another_input && !helper.IsEmpty(validInput) {
			all_absent = false
			break
		}
	}

	if _, exists_input := payload[input]; !exists_input && all_absent {
		all_inputs := strings.Join(options, ", ")
		tmpError := fmt.Sprintf("El campo \"%s\" debe estar definido cuando los campos \"%s\" no están definidos o están vacíos", input, all_inputs)

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo \"%s\" en la posición %s debe estar definido cuando los campos \"%s\" no están definidos o están vacíos", input, sliceIndex, all_inputs)
		}

		tmpErrorKey := fmt.Sprintf("%s.required_without_all", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "required_without_all", errors, tmpError)
		return errors, true
	}

	return errors, false
}
