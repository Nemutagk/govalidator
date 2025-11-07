package validate

import (
	"fmt"
	"strconv"
	"time"
)

func Before(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if len(options) == 0 {
		errors = addError(input, "before", errors, "El valor a comparar no está definido")
		return errors
	}

	value := payload[input]

	if value == nil || value == "" {
		return errors
	}

	date, err_date := time.Parse("2006-01-02", value.(string))

	if err_date == nil {
		if options[0] == "now" || options[0] == "current" || options[0] == "today" {
			if date.After(time.Now()) {
				tmpError := "La fecha no es anterior a la fecha actual"

				customeErrorKey := fmt.Sprintf("%s.before", input)
				if customeError, exists := customeErrors[customeErrorKey]; exists {
					tmpError = customeError
				}

				errors = addError(input, "before", errors, tmpError)
			}
		}

		if options[0] == "tomorrow" {
			if date.After(time.Now().AddDate(0, 0, 1)) {
				tmpError := "La fecha no es anterior a mañana"

				customeErrorKey := fmt.Sprintf("%s.before", input)
				if customeError, exists := customeErrors[customeErrorKey]; exists {
					tmpError = customeError
				}

				errors = addError(input, "before", errors, tmpError)
			}
		}

		if options[0] == "yesterday" {
			if date.After(time.Now().AddDate(0, 0, -1)) {
				tmpError := "La fecha no es anterior a ayer"

				customeErrorKey := fmt.Sprintf("%s.before", input)
				if customeError, exists := customeErrors[customeErrorKey]; exists {
					tmpError = customeError
				}

				errors = addError(input, "before", errors, tmpError)
			}
		}

		compare_date, err := time.Parse("2006-01-02", options[0])
		if err == nil {
			if date.After(compare_date) {
				tmpError := "La fecha no es anterior a la fecha " + options[0]

				customeErrorKey := fmt.Sprintf("%s.before", input)
				if customeError, exists := customeErrors[customeErrorKey]; exists {
					tmpError = customeError
				}

				errors = addError(input, "before", errors, tmpError)
			}
		}
	}

	num, ok := value.(int)
	if ok {
		if len(options) == 0 {
			errors = addError(input, "before", errors, "El valor a comparar no está definido")
		} else {
			compare_int, _ := strconv.Atoi(options[0])
			if num > compare_int {
				tmpError := "La entrada " + input + " no es anterior al número " + options[0]

				customeErrorKey := fmt.Sprintf("%s.before", input)
				if customeError, exists := customeErrors[customeErrorKey]; exists {
					tmpError = customeError
				}

				errors = addError(input, "before", errors, tmpError)
			}
		}
	}

	return errors
}
