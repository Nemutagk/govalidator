package validate

import "reflect"

func Array(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	value, exsits := payload[input]
	if !exsits {
		return errors
	}

	v := reflect.ValueOf(value)

	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		addError(input, "before", errors, "El valor no es un array.")
		return errors
	}

	return errors
}
