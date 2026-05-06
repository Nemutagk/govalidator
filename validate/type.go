package validate

import (
	"fmt"
	"reflect"
	"slices"
)

func Type(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	value, exist := payload[input]

	if !exist {
		return errors
	}

	if len(options) == 0 {
		tmpError := "El tipo no está definido"

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El tipo en la posición %s no está definido", sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.type", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "type", errors, tmpError)
		return errors
	}

	if value == nil {
		if slices.Contains(options, "nullable") {
			return errors
		}

		tmpError := fmt.Sprintf("El campo \"%s\" no puede ser nulo", input)

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo \"%s\" en la posición %s no puede ser nulo", input, sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.type", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "type", errors, tmpError)
		return errors
	}

	var_type := reflect.TypeOf(value).String()
	if var_type != options[0] {
		tmpError := fmt.Sprintf("El tipo del campo \"%s\" no es \"%s\"", input, options[0])

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El tipo en la posición %s del campo \"%s\" no es \"%s\"", sliceIndex, input, options[0])
		}

		tmpErrorKey := fmt.Sprintf("%s.type", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "type", errors, tmpError)
		return errors
	}

	return errors
}
