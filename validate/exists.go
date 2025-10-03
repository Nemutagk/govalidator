package validate

func Exists(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, modelList map[string]func(data string) bool) map[string]interface{} {
	if len(options) != 1 {
		errors = addError(input, "exists", errors, "la configuración de conexión no es válida")
		return errors
	}

	model, ok := modelList[options[0]]
	if !ok || model == nil {
		errors = addError(input, "exists", errors, "el modelo no está definido")
		return errors
	}

	value, ok := payload[input].(string)
	if !ok {
		errors = addError(input, "exists", errors, "el valor no es válido")
		return errors
	}

	result := model(value)

	if !result {
		errors = addError(input, "exists", errors, "El valor '"+value+"' no existe")
	}

	return errors
}
