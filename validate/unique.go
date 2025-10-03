package validate

import (
	"fmt"
)

/**
* Validates that the indicated value does not exists in the database
* Params: db connection, table name, column name
* Example: unique:princial,users,email
 */
func Unique(input string, payload map[string]interface{}, options []string, list_errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, listModels map[string]func(data string) bool) map[string]interface{} {
	if len(options) != 1 {
		list_errors = addError(input, "unique", list_errors, "the options is not valid")
		return list_errors
	}

	if _, exists_input := payload[input]; !exists_input {
		fmt.Println("validate unique:input not exists")
		return list_errors
	}

	if listModels == nil {
		list_errors = addError(input, "unique", list_errors, "the model is not defined")
		return list_errors
	}

	model, ok := listModels[options[0]]
	if !ok || model == nil {
		list_errors = addError(input, "unique", list_errors, "the model is not defined")
		return list_errors
	}

	value, ok := payload[input].(string)
	if !ok {
		list_errors = addError(input, "unique", list_errors, "el valor no es válido")
		return list_errors
	}

	result := model(value)

	if !result {
		list_errors = addError(input, "unique", list_errors, "El valor '"+value+"'ya está registrado")
	}

	return list_errors
}
