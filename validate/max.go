package validate

import (
	"fmt"
	"reflect"
	"strconv"
)

func Max(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	max, err := strconv.ParseInt(options[0], 10, 64)
	if err != nil {
		tmpError := fmt.Sprintf("El campo %s debe ser un número", input)

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo %s en la posición %s debe ser un número", input, sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.max", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "max", errors, tmpError)
		return errors
	}

	if _, ok := value.(string); ok {
		strlen := len(value.(string))

		if strlen > int(max) {
			tmpError := fmt.Sprintf("El campo %s debe tener como máximo %s caracteres", input, options[0])

			if sliceIndex != "" {
				tmpError = fmt.Sprintf("El campo %s en la posición %s debe tener como máximo %s caracteres", input, sliceIndex, options[0])
			}

			tmpErrorKey := fmt.Sprintf("%s.max", input)
			if customeError, exists := customeErrors[tmpErrorKey]; exists {
				tmpError = customeError
			}
			errors = addError(input, "max", errors, tmpError)
		}
	}

	if _, ok := value.(int); ok {
		intValue := value.(int)
		if intValue > int(max) {
			tmpError := fmt.Sprintf("El campo %s debe ser como máximo %s", input, options[0])

			if sliceIndex != "" {
				tmpError = fmt.Sprintf("El campo %s en la posición %s debe ser como máximo %s", input, sliceIndex, options[0])
			}

			tmpErrorKey := fmt.Sprintf("%s.max", input)
			if customeError, exists := customeErrors[tmpErrorKey]; exists {
				tmpError = customeError
			}
			errors = addError(input, "max", errors, tmpError)
		}
	}

	if _, ok := value.(float64); ok {
		floatValue := value.(float64)
		if floatValue > float64(max) {
			tmpError := fmt.Sprintf("El campo %s debe ser como máximo %s", input, options[0])
			tmpErrorKey := fmt.Sprintf("%s.max", input)
			if customeError, exists := customeErrors[tmpErrorKey]; exists {
				tmpError = customeError
			}
			errors = addError(input, "max", errors, tmpError)
		}
	}

	val := reflect.ValueOf(value)
	kind := val.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		if val.Len() > int(max) {
			tmpError := fmt.Sprintf("El campo %s debe tener como máximo %s elementos", input, options[0])

			if sliceIndex != "" {
				tmpError = fmt.Sprintf("El campo %s en la posición %s debe tener como máximo %s elementos", input, sliceIndex, options[0])
			}

			tmpErrorKey := fmt.Sprintf("%s.max", input)
			if customeError, exists := customeErrors[tmpErrorKey]; exists {
				tmpError = customeError
			}
			errors = addError(input, "max", errors, tmpError)
		}
	}

	return errors
}
