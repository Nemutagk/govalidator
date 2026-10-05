package validate

import (
	"fmt"
	"reflect"
	"strconv"
	"unicode/utf8"
)

func Len(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if len(options) == 0 {
		tmpError := fmt.Sprintf("La regla len del campo %s requiere un valor", input)

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("La regla len del campo %s en la posición %s requiere un valor", input, sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.len", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "len", errors, tmpError)
		return errors
	}

	expected, err := strconv.ParseInt(options[0], 10, 64)
	if err != nil {
		tmpError := fmt.Sprintf("El campo %s debe ser un número", input)

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo %s en la posición %s debe ser un número", input, sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.len", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "len", errors, tmpError)
		return errors
	}

	val := reflect.ValueOf(value)

	var total int
	var unit string
	switch val.Kind() {
	case reflect.String:
		total = utf8.RuneCountInString(val.String())
		unit = "caracteres"
	case reflect.Slice, reflect.Array, reflect.Map:
		total = val.Len()
		unit = "elementos"
	default:
		return errors
	}

	if int64(total) != expected {
		tmpError := fmt.Sprintf("El campo %s debe tener exactamente %s %s", input, options[0], unit)

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo %s en la posición %s debe tener exactamente %s %s", input, sliceIndex, options[0], unit)
		}

		tmpErrorKey := fmt.Sprintf("%s.len", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "len", errors, tmpError)
	}

	return errors
}
