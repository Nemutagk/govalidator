package validate

import (
	"fmt"
	"reflect"
	"strconv"
)

func Min(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	min, err := strconv.ParseInt(options[0], 10, 64)
	if err != nil {
		tmpError := fmt.Sprintf("El campo %s debe ser un número", input)

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo %s en la posición %s debe ser un número", input, sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.min", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "min", errors, tmpError)
		return errors
	}

	if _, ok := value.(string); ok {
		strlen := len(value.(string))

		if strlen < int(min) {
			tmpError := fmt.Sprintf("El campo %s debe tener al menos %s caracteres", input, options[0])

			if sliceIndex != "" {
				tmpError = fmt.Sprintf("El campo %s en la posición %s debe tener al menos %s caracteres", input, sliceIndex, options[0])
			}

			tmpErrorKey := fmt.Sprintf("%s.min", input)
			if customeError, exists := customeErrors[tmpErrorKey]; exists {
				tmpError = customeError
			}
			errors = addError(input, "min", errors, tmpError)
		}
	}

	if _, ok := value.(int); ok {
		intValue := value.(int)
		if intValue < int(min) {
			tmpError := fmt.Sprintf("El campo %s debe ser al menos %s", input, options[0])

			if sliceIndex != "" {
				tmpError = fmt.Sprintf("El campo %s en la posición %s debe ser al menos %s", input, sliceIndex, options[0])
			}

			tmpErrorKey := fmt.Sprintf("%s.min", input)
			if customeError, exists := customeErrors[tmpErrorKey]; exists {
				tmpError = customeError
			}
			errors = addError(input, "min", errors, tmpError)
		}
	}

	if _, ok := value.(float64); ok {
		floatValue := value.(float64)
		if floatValue < float64(min) {
			tmpError := fmt.Sprintf("El campo %s debe ser al menos %s", input, options[0])

			if sliceIndex != "" {
				tmpError = fmt.Sprintf("El campo %s en la posición %s debe ser al menos %s", input, sliceIndex, options[0])
			}

			tmpErrorKey := fmt.Sprintf("%s.min", input)
			if customeError, exists := customeErrors[tmpErrorKey]; exists {
				tmpError = customeError
			}
			errors = addError(input, "min", errors, tmpError)
		}
	}

	val := reflect.ValueOf(value)
	kind := val.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		if val.Len() < int(min) {
			tmpError := fmt.Sprintf("El campo %s debe tener al menos %s elementos", input, options[0])

			if sliceIndex != "" {
				tmpError = fmt.Sprintf("El campo %s en la posición %s debe tener al menos %s elementos", input, sliceIndex, options[0])
			}

			tmpErrorKey := fmt.Sprintf("%s.min", input)
			if customeError, exists := customeErrors[tmpErrorKey]; exists {
				tmpError = customeError
			}
			errors = addError(input, "min", errors, tmpError)
		}
	}

	return errors
}
