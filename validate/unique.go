package validate

import (
	"fmt"
)

/**
* Validates that the indicated value does not exists in the database
* Params: db connection, table name, column name
* Example: unique:princial,users,email
 */
func Unique(input string, value any, payload map[string]any, options []string, sliceIndex string, list_errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, listModels map[string]func(data string, payload map[string]any) bool, customeErrors map[string]string) map[string]interface{} {
	if len(options) != 1 {
		list_errors = addError(input, "unique", list_errors, "the options is not valid")
		return list_errors
	}

	if _, exists_input := payload[input]; !exists_input {
		fmt.Println("validate unique:input not exists")
		return list_errors
	}

	if listModels == nil {
		if sliceIndex != "" {
			list_errors = addError(input, "unique", list_errors, fmt.Sprintf("el valor en la posición %s no es válido", sliceIndex))
		} else {
			list_errors = addError(input, "unique", list_errors, "el valor no es válido")
		}
		return list_errors
	}

	model, ok := listModels[options[0]]
	if !ok || model == nil {
		if sliceIndex != "" {
			list_errors = addError(input, "unique", list_errors, fmt.Sprintf("el valor en la posición %s no es válido", sliceIndex))
		} else {
			list_errors = addError(input, "unique", list_errors, "el valor no es válido")
		}
		return list_errors
	}

	valueStr, ok := value.(string)
	if !ok {
		if sliceIndex != "" {
			list_errors = addError(input, "unique", list_errors, fmt.Sprintf("el valor en la posición %s no es válido", sliceIndex))
		} else {
			list_errors = addError(input, "unique", list_errors, "el valor no es válido")
		}
		return list_errors
	}

	result := model(valueStr, payload)

	if !result {
		tmpError := "El valor '" + valueStr + "' ya está registrado"

		customeErrorKey := fmt.Sprintf("%s.unique", input)
		if customeError, exists := customeErrors[customeErrorKey]; exists {
			tmpError = customeError
		}

		list_errors = addError(input, "unique", list_errors, tmpError)
	}

	return list_errors
}
