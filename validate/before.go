package validate

import (
	"fmt"
	"time"
)

func Before(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if len(options) == 0 {
		errors = addError(input, "before", errors, "El valor a comparar no está definido")
		return errors
	}

	if value == nil || value == "" {
		return errors
	}

	if ownNum, isNumber := toFloat64(value); isNumber {
		cmpNum, isField, errMsg := resolveNumericTarget(payload, options[0])
		if errMsg != "" {
			errors = addError(input, "before", errors, errMsg)
			return errors
		}

		if ownNum > cmpNum {
			targetLabel := "número"
			if isField {
				targetLabel = "campo"
			}

			tmpError := fmt.Sprintf("La entrada %s no es anterior al %s %s", input, targetLabel, options[0])

			if sliceIndex != "" {
				tmpError = fmt.Sprintf("La entrada en la posición %s no es anterior al %s %s", sliceIndex, targetLabel, options[0])
			}

			customeErrorKey := fmt.Sprintf("%s.before", input)
			if customeError, exists := customeErrors[customeErrorKey]; exists {
				tmpError = customeError
			}

			errors = addError(input, "before", errors, tmpError)
		}

		return errors
	}

	formato := "2006-01-02"
	if len(options) >= 2 {
		// log.Printf("Formato personalizado detectado: %s", options[1])
		formato = options[1]
	}

	valueStr, isString := value.(string)
	if !isString {
		errors = addError(input, "before", errors, "La fecha proporcionada no es válida o no coincide con el formato "+formato)
		return errors
	}

	date, err_date := time.Parse(formato, valueStr)
	if err_date != nil {
		errors = addError(input, "before", errors, "La fecha proporcionada no es válida o no coincide con el formato "+formato)
		return errors
	}

	if len(options) > 0 {
		if fecha_str, ok := payload[options[0]]; ok {
			formato = "2006-01-02"
			if len(options) >= 2 {
				formato = options[1]
			}
			fecha_str_val, isCompareString := fecha_str.(string)
			if !isCompareString {
				errors = addError(input, "before", errors, "La fecha de comparación no es válida o no coincide con el formato "+formato)
				return errors
			}

			fecha_comparar, err_fecha := time.Parse(formato, fecha_str_val)

			if err_fecha != nil {
				errors = addError(input, "before", errors, "La fecha de comparación no es válida o no coincide con el formato "+formato)
				return errors
			}
			if !date.Before(fecha_comparar) {
				tmpError := "La fecha no es anterior a la fecha " + fecha_str_val

				customeErrorKey := fmt.Sprintf("%s.before", input)
				if customeError, exists := customeErrors[customeErrorKey]; exists {
					tmpError = customeError
				}

				errors = addError(input, "before", errors, tmpError)
			}
			return errors
		}
	}

	humanDays := map[string]string{
		"now":       "now",
		"current":   "now",
		"today":     "now",
		"tomorrow":  "tomorrow",
		"yesterday": "yesterday",
	}

	if _, exists := humanDays[options[0]]; exists {
		if options[0] == "now" || options[0] == "current" || options[0] == "today" {
			// log.Printf("se compara con now")
			if !date.Before(time.Now()) {
				tmpError := "La fecha no es anterior a la fecha actual"

				customeErrorKey := fmt.Sprintf("%s.before", input)
				if customeError, exists := customeErrors[customeErrorKey]; exists {
					tmpError = customeError
				}

				errors = addError(input, "before", errors, tmpError)
			}
		}

		if options[0] == "tomorrow" {
			if !date.Before(time.Now().AddDate(0, 0, 1)) {
				tmpError := "La fecha no es anterior a mañana"

				customeErrorKey := fmt.Sprintf("%s.before", input)
				if customeError, exists := customeErrors[customeErrorKey]; exists {
					tmpError = customeError
				}

				errors = addError(input, "before", errors, tmpError)
			}
		}

		if options[0] == "yesterday" {
			if !date.Before(time.Now().AddDate(0, 0, -1)) {
				tmpError := "La fecha no es anterior a ayer"

				customeErrorKey := fmt.Sprintf("%s.before", input)
				if customeError, exists := customeErrors[customeErrorKey]; exists {
					tmpError = customeError
				}

				errors = addError(input, "before", errors, tmpError)
			}
		}

		return errors
	}

	compare_date, err := time.Parse(formato, options[0])
	if err != nil {
		errors = addError(input, "before", errors, "La fecha de comparación no es válida o no coincide con el formato "+formato)
		return errors
	}

	if !date.Before(compare_date) {
		tmpError := "La fecha no es anterior a la fecha " + options[0]

		customeErrorKey := fmt.Sprintf("%s.before", input)
		if customeError, exists := customeErrors[customeErrorKey]; exists {
			tmpError = customeError
		}

		errors = addError(input, "before", errors, tmpError)
	}

	return errors
}
