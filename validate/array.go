package validate

import (
	"fmt"
	"reflect"
)

func Array(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	v := reflect.ValueOf(value)

	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		tmpError := "El valor no es un array."

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo %s en la posición %s no es un array.", input, sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.array", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		addError(input, "before", errors, tmpError)
		return errors
	}

	return errors
}
