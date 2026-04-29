package validate

import "fmt"

func Customized(input string, value any, payload map[string]any, options []string, sliceIndex string, list_errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, listModels map[string]func(data any, payload map[string]any, opts *[]string) (bool, string), customizedErrors map[string]string) map[string]interface{} {
	if len(options) == 0 {
		if sliceIndex != "" {
			list_errors = addError(input, "customized", list_errors, fmt.Sprintf("La función de validación personalizada en la posición %s no está definida", sliceIndex))
		} else {
			list_errors = addError(input, "customized", list_errors, "La función de validación personalizada no está definida")
		}
		return list_errors
	}

	customizedFunc, exists := listModels[options[0]]
	if !exists {
		if sliceIndex != "" {
			list_errors = addError(input, "customized", list_errors, fmt.Sprintf("La función de validación personalizada en la posición %s no existe", sliceIndex))
		} else {
			list_errors = addError(input, "customized", list_errors, "La función de validación personalizada no existe")
		}
		return list_errors
	}

	lastOpts := options[1:]
	result, customErr := customizedFunc(value, payload, &lastOpts)

	if !result {
		tmpError := "La validación personalizada ha fallado"

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("La validación personalizada en la posición %s ha fallado", sliceIndex)
		}

		customizedErrorKey := fmt.Sprintf("%s.customized", input)
		if customizedError, exists := customizedErrors[customizedErrorKey]; exists {
			tmpError = customizedError
		}

		if customErr != "" {
			tmpError = customErr
		}

		list_errors = addError(input, "customized", list_errors, tmpError)
	}

	return list_errors
}
