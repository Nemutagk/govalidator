package validate

import (
	"fmt"
	"reflect"
)

func Type(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	value, exist := payload[input]

	if !exist {
		return errors
	}

	if len(options) == 0 {
		errors = addError(input, "type", errors, "El tipo no está definido")
		return errors
	}

	var_type := reflect.TypeOf(value).String()
	fmt.Println("var_type: ", var_type)
	if var_type != options[0] {
		errors = addError(input, "type", errors, "El tipo del campo "+input+" no es "+options[0])
		return errors
	}

	return errors
}
