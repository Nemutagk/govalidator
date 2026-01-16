package validate

import "fmt"

func Exists(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, modelList map[string]func(data any, payload map[string]any, opts *[]string) bool, customeErrors map[string]string) map[string]interface{} {
	if len(options) != 1 {
		tmpError := "la configuración de conexión no es válida"

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("La configuración de conexión en la posición %s no es válida", sliceIndex)
		}

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

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El modelo en la posición %s no está definido", sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.exists", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "exists", errors, tmpError)
		return errors
	}

	valueStr, ok := value.(string)
	if !ok {
		tmpError := "el valor no es válido"

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El valor en la posición %s no es válido", sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.exists", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "exists", errors, tmpError)
		return errors
	}

	anotherOpts := options[1:]
	result := model(valueStr, payload, &anotherOpts)

	if !result {
		tmpError := fmt.Sprintf("El valor '%s' no existe", value)

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El valor '%s' en la posición %s no existe", value, sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.exists", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "exists", errors, tmpError)
	}

	return errors
}
