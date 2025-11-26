package validate

import (
	"fmt"
	"strconv"
	"time"
)

func After(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if len(options) == 0 {
		errors = addError(input, "after", errors, "El valor a comprar no esta definido")
		return errors
	}

	if value == nil || value == "" {
		return errors
	}

	num, ok := value.(int)
	if ok {
		if len(options) == 0 {
			if sliceIndex != "" {
				input = fmt.Sprintf("%s[%s]", input, sliceIndex)
			} else {
				errors = addError(input, "after", errors, "El valor a comparar no está definido")
			}
		} else {
			compare_int, _ := strconv.Atoi(options[0])
			if num < compare_int {
				tmpError := "La entrada " + input + " no es posterior al número " + options[0]

				if sliceIndex != "" {
					tmpError = fmt.Sprintf("La entrada en la posición %s no es posterior al número %s", sliceIndex, options[0])
				}

				customeErrorKey := fmt.Sprintf("%s.after", input)
				if customeError, exists := customeErrors[customeErrorKey]; exists {
					tmpError = customeError
				}

				errors = addError(input, "after", errors, tmpError)
			}
		}
	}

	formato := "2006-01-02"
	if len(options) >= 2 {
		formato = options[1]
	}

	date, err_date := time.Parse(formato, value.(string))
	if err_date != nil {
		errors = addError(input, "after", errors, "El valor no es una fecha válida")
		return errors
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
			if !date.After(time.Now()) {
				tmpError := "La fecha no es posterior a la fecha actual"

				customeErrorKey := fmt.Sprintf("%s.after", input)
				if customeError, exists := customeErrors[customeErrorKey]; exists {
					tmpError = customeError
				}

				errors = addError(input, "after", errors, tmpError)
			}
		}

		if options[0] == "tomorrow" {
			if !date.After(time.Now().AddDate(0, 0, 1)) {
				tmpError := "La fecha no es posterior a mañana"

				customeErrorKey := fmt.Sprintf("%s.after", input)
				if customeError, exists := customeErrors[customeErrorKey]; exists {
					tmpError = customeError
				}

				errors = addError(input, "after", errors, tmpError)
			}
		}

		if options[0] == "yesterday" {
			if !date.After(time.Now().AddDate(0, 0, -1)) {
				tmpError := "La fecha no es posterior a ayer"

				customeErrorKey := fmt.Sprintf("%s.after", input)
				if customeError, exists := customeErrors[customeErrorKey]; exists {
					tmpError = customeError
				}

				errors = addError(input, "after", errors, tmpError)
			}
		}

		return errors
	}

	compare_date, err := time.Parse(formato, options[0])
	if err == nil {
		if !date.After(compare_date) {
			tmpError := "La fecha no es posterior a la fecha " + options[0]

			customeErrorKey := fmt.Sprintf("%s.after", input)
			if customeError, exists := customeErrors[customeErrorKey]; exists {
				tmpError = customeError
			}

			errors = addError(input, "after", errors, tmpError)
		}
	}

	return errors
}
