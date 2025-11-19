package validate

import (
	"fmt"
	"strconv"
	"time"
)

func After(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if len(options) == 0 {
		errors = addError(input, "after", errors, "El valor a comprar no esta definido")
		return errors
	}

	value := payload[input]

	if value == nil || value == "" {
		return errors
	}

	formato := "2006-01-02T15:04:05"
	if len(options) >= 2 {
		formato = options[1]
	}

	date, err_date := time.Parse(formato, value.(string))

	if err_date == nil {
		if options[0] == "now" || options[0] == "current" || options[0] == "today" {
			if date.Before(time.Now()) {
				tmpError := "La fecha no es posterior a la fecha actual"

				customeErrorKey := fmt.Sprintf("%s.after", input)
				if customeError, exists := customeErrors[customeErrorKey]; exists {
					tmpError = customeError
				}

				errors = addError(input, "after", errors, tmpError)
			}
		}

		if options[0] == "tomorrow" {
			if date.Before(time.Now().AddDate(0, 0, 1)) {
				tmpError := "La fecha no es posterior a mañana"

				customeErrorKey := fmt.Sprintf("%s.after", input)
				if customeError, exists := customeErrors[customeErrorKey]; exists {
					tmpError = customeError
				}

				errors = addError(input, "after", errors, tmpError)
			}
		}

		if options[0] == "yesterday" {
			if date.Before(time.Now().AddDate(0, 0, -1)) {
				tmpError := "La fecha no es posterior a ayer"

				customeErrorKey := fmt.Sprintf("%s.after", input)
				if customeError, exists := customeErrors[customeErrorKey]; exists {
					tmpError = customeError
				}

				errors = addError(input, "after", errors, tmpError)
			}
		}

		compare_date, err := time.Parse(formato, options[0])
		if err == nil {
			if date.Before(compare_date) {
				tmpError := "La fecha no es posterior a la fecha " + options[0]

				customeErrorKey := fmt.Sprintf("%s.after", input)
				if customeError, exists := customeErrors[customeErrorKey]; exists {
					tmpError = customeError
				}

				errors = addError(input, "after", errors, tmpError)
			}
		}
	}

	num, ok := value.(int)
	if ok {
		if len(options) == 0 {
			errors = addError(input, "after", errors, "El valor a comparar no está definido")
		} else {
			compare_int, _ := strconv.Atoi(options[0])
			if num < compare_int {
				tmpError := "La entrada " + input + " no es posterior al número " + options[0]

				customeErrorKey := fmt.Sprintf("%s.after", input)
				if customeError, exists := customeErrors[customeErrorKey]; exists {
					tmpError = customeError
				}

				errors = addError(input, "after", errors, tmpError)
			}
		}
	}

	return errors
}
