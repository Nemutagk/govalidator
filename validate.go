package govalidator

import (
	"github.com/Nemutagk/goerrors"
	"github.com/Nemutagk/govalidator/v2/validate"
)

type Input struct {
	Name  string
	Rules []Rule
}

type Rule struct {
	Name    string
	Options []string
}

func ValidateRequest(body map[string]any, rules []Input, customeErrors map[string]string, models map[string]func(data string, payload map[string]any) bool) (map[string]any, *goerrors.GError) {
	errors := make(map[string]any)
	safePayload := make(map[string]any)

	for _, input := range rules {
		skipRulesMap := false
		inputName := input.Name

		for _, rule := range input.Rules {
			opts := rule.Options

			switch rule.Name {
			case "email":
				errors = validate.Email(inputName, body, opts, errors, addError, customeErrors)
			case "confirmation":
				errors = validate.Confirmation(inputName, body, opts, errors, addError, customeErrors)
			case "unique":
				errors = validate.Unique(inputName, body, opts, errors, addError, models, customeErrors)
			case "in":
				errors = validate.In(inputName, body, opts, errors, addError, customeErrors)
			case "not_in":
				errors = validate.NotIn(inputName, body, opts, errors, addError, customeErrors)
			case "before":
				errors = validate.Before(inputName, body, opts, errors, addError, customeErrors)
			case "after":
				errors = validate.After(inputName, body, opts, errors, addError, customeErrors)
			case "ip":
				errors = validate.Ip(inputName, body, opts, errors, addError, customeErrors)
			case "password":
				errors = validate.Password(inputName, body, opts, errors, addError, customeErrors)
			case "exists":
				errors = validate.Exists(inputName, body, opts, errors, addError, models, customeErrors)
			case "min":
				errors = validate.Min(inputName, body, opts, errors, addError, customeErrors)
			case "max":
				errors = validate.Max(inputName, body, opts, errors, addError, customeErrors)
			case "boolean":
				errors = validate.Boolean(inputName, body, opts, errors, addError, customeErrors)
			case "sometimes":
				if _, exists_input := body[inputName]; !exists_input {
					//Si no existe el input no se validan lás demás reglas existentes
					// log.Println("Input", input, "no existe en el body, no se validan las demás reglas")
					skipRulesMap = true
				}
			case "required":
				errors = validate.Required(inputName, body, opts, errors, addError, customeErrors)
			case "required_with":
				errors = validate.RequiredWith(inputName, body, opts, errors, addError, customeErrors)
			case "required_with_all":
				errors = validate.RequiredWithAll(inputName, body, opts, errors, addError, customeErrors)
			case "required_without":
				errors = validate.RequiredWithout(inputName, body, opts, errors, addError, customeErrors)
			case "required_without_all":
				errors = validate.RequiredWithoutAll(inputName, body, opts, errors, addError, customeErrors)
			case "array":
				errors = validate.Array(inputName, body, opts, errors, addError, customeErrors)
			case "type":
				errors = validate.Type(inputName, body, opts, errors, addError, customeErrors)
			case "date":
				errors = validate.Date(inputName, body, opts, errors, addError, customeErrors)
			case "date_format":
				errors = validate.DateFormat(inputName, body, opts, errors, addError, customeErrors)
			case "custome":
				errors = validate.Custome(inputName, body, opts, errors, addError, models, customeErrors)

			default:
				errors = addError(inputName, rule.Name, errors, "The rule "+rule.Name+" is not valid")
			}

			if skipRulesMap {
				break
			}
		}

		if _, ok := body[inputName]; ok {
			safePayload[inputName] = body[inputName]
		}
	}

	allErrors := make([]string, 0)
	if len(errors) > 0 {
		for input, inputErrors := range errors {
			for _, errMessages := range inputErrors.(map[string]interface{}) {
				for _, errMessage := range errMessages.([]string) {
					allErrors = append(allErrors, input+": "+errMessage)
				}
			}
		}
		gerr := goerrors.NewGError("Error en la validación", goerrors.StatusBadRequest, &allErrors, nil)

		return nil, gerr
	}

	return safePayload, nil
}

func addError(input string, rule string, errors map[string]interface{}, error string) map[string]interface{} {
	if _, exists_input := errors[input]; !exists_input {
		errors[input] = map[string]interface{}{
			rule: []string{
				error,
			},
		}
	} else {
		if inputErrors, ok := errors[input].(map[string]interface{}); ok {
			if _, exists_rule := inputErrors[rule]; !exists_rule {
				inputErrors[rule] = []string{
					error,
				}
			} else {
				inputErrors[rule] = append(inputErrors[rule].([]string), error)
			}
			errors[input] = inputErrors
		} else {
			errors[input] = map[string]interface{}{
				rule: []string{
					error,
				},
			}
		}
	}

	return errors
}
