package validate

import "fmt"

func Exists(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, modelList map[string]func(data string, payload map[string]any) bool, customeErrors map[string]string) map[string]interface{} {
	if len(options) != 1 {
		tmpError := "la configuración de conexión no es válida"
		tmpErrorKey := fmt.Sprintf("%s.exists", input)
		if _, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeErrors[tmpErrorKey]
		}
		errors = addError(input, "exists", errors, tmpError)
		return errors
	}

	model, ok := modelList[options[0]]
	if !ok || model == nil {
		tmpError := "el modelo no está definido"
		tmpErrorKey := fmt.Sprintf("%s.exists", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "exists", errors, tmpError)
		return errors
	}

	value, ok := payload[input].(string)
	if !ok {
		tmpError := "el valor no es válido"
		tmpErrorKey := fmt.Sprintf("%s.exists", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "exists", errors, tmpError)
		return errors
	}

	result := model(value, payload)

	if !result {
		tmpError := fmt.Sprintf("El valor '%s' no existe", value)
		tmpErrorKey := fmt.Sprintf("%s.exists", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "exists", errors, tmpError)
	}

	return errors
}
