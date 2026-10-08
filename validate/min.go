package validate

import (
	"fmt"
	"reflect"
	"strconv"
	"unicode/utf8"
)

func Min(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if len(options) == 0 {
		tmpError := fmt.Sprintf("La regla min del campo %s requiere un valor", input)

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("La regla min del campo %s en la posición %s requiere un valor", input, sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.min", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "min", errors, tmpError)
		return errors
	}

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

	if strValue, ok := value.(string); ok {
		strlen := utf8.RuneCountInString(strValue)

		if int64(strlen) < min {
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

	if result, isNumber := compareToInt64(value, min); isNumber && result < 0 {
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

	val := reflect.ValueOf(value)
	kind := val.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		if int64(val.Len()) < min {
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
