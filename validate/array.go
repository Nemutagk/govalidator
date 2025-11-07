package validate

import (
	"fmt"
	"reflect"
)

func Array(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	value, exsits := payload[input]
	if !exsits {
		return errors
	}

	v := reflect.ValueOf(value)

	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		tmpError := "El valor no es un array."
		tmpErrorKey := fmt.Sprintf("%s.array", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		addError(input, "before", errors, tmpError)
		return errors
	}

	return errors
}
