package validate

import "fmt"

func Custome(input string, payload map[string]interface{}, options []string, list_errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, listModels map[string]func(data string, payload map[string]any) bool, customeErrors map[string]string) map[string]interface{} {
	if len(options) == 0 {
		list_errors = addError(input, "custome", list_errors, "La función de validación personalizada no está definida")
		return list_errors
	}

	customeFunc, exists := listModels[options[0]]
	if !exists {
		list_errors = addError(input, "custome", list_errors, "La función de validación personalizada no existe")
		return list_errors
	}

	value := payload[input].(string)
	result := customeFunc(value, payload)

	if !result {
		tmpError := "La validación personalizada ha fallado"

		customeErrorKey := fmt.Sprintf("%s.custome", input)
		if customeError, exists := customeErrors[customeErrorKey]; exists {
			tmpError = customeError
		}

		list_errors = addError(input, "custome", list_errors, tmpError)
	}

	return list_errors
}
