package validate

import (
	"fmt"
	"reflect"
)

func Type(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	value, exist := payload[input]

	if !exist {
		return errors
	}

	if len(options) == 0 {
		tmpError := "El tipo no está definido"
		tmpErrorKey := fmt.Sprintf("%s.type", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "type", errors, tmpError)
		return errors
	}

	var_type := reflect.TypeOf(value).String()
	fmt.Println("var_type: ", var_type)
	if var_type != options[0] {
		tmpError := fmt.Sprintf("El tipo del campo \"%s\" no es \"%s\"", input, options[0])
		tmpErrorKey := fmt.Sprintf("%s.type", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "type", errors, tmpError)
		return errors
	}

	return errors
}
