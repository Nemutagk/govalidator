package govalidator

import (
	"encoding/json"
	"strings"

	"github.com/Nemutagk/godb"
	"github.com/Nemutagk/godb/definitions/db"
	"github.com/Nemutagk/goerrors"
	"github.com/Nemutagk/govalidator/validate"
)

func ValidateRequestT[T any](body T, rules map[string]string) (*T, *goerrors.GError) {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, goerrors.NewGError("Error en la validación", goerrors.StatusBadRequest, &[]string{
			"Error al codificar el payload: " + err.Error(),
		}, nil)
	}
	var bodyMap map[string]interface{}
	err = json.Unmarshal(bodyJson, &bodyMap)
	if err != nil {
		return nil, goerrors.NewGError("Error en la validación", goerrors.StatusBadRequest, &[]string{
			"Error al decodificar el payload: " + err.Error(),
		}, nil)
	}

	var safeResult T
	errMap := ValidateRequest(bodyMap, rules, &safeResult)
	if errMap != nil {
		return nil, errMap
	}

	return &safeResult, nil
}

func ValidateRequest(body map[string]interface{}, rules map[string]string, result any) *goerrors.GError {

	rules_inputs := make(map[string]interface{})

	body_encode, err := json.Marshal(body)

	if err != nil {
		convertError := goerrors.ConvertError(err)
		return goerrors.NewGError("Error en la validación", goerrors.StatusBadRequest, &[]string{
			"The payload is invalid",
		}, convertError)
	}

	var body_parse map[string]interface{}

	err = json.Unmarshal(body_encode, &body_parse)

	if err != nil {
		convertError := goerrors.ConvertError(err)
		return goerrors.NewGError("Error en la validación", goerrors.StatusBadRequest, &[]string{
			"The payload is invalid",
		}, convertError)
	}

	for input, rules := range rules {
		separated_rules := strings.Split(rules, "|")
		rules_building := make(map[string]interface{})
		for _, rule := range separated_rules {
			tmp_rule := getRuleWithOptions(rule)
			rules_building[tmp_rule[0]] = tmp_rule[1:]
		}
		rules_inputs[input] = rules_building
	}

	errors := make(map[string]interface{})
	safePayload := make(map[string]interface{})

	dbManager := godb.InitConnections(map[string]db.DbConnection{})

	for input, rules := range rules_inputs {
		rulesMap, ok := rules.(map[string]interface{})
		if !ok {
			continue
		}

		skipRulesMap := false

		for rule, options := range rulesMap {
			opts, ok := options.([]string)
			if !ok {
				continue
			}

			switch rule {
			case "email":
				errors = validate.Email(input, body_parse, opts, errors, addError)
			case "confirmation":
				errors = validate.Confirmation(input, body_parse, opts, errors, addError)
			case "unique":
				errors = validate.Unique(input, body_parse, opts, errors, addError, dbManager)
			case "in":
				errors = validate.In(input, body_parse, opts, errors, addError)
			case "before":
				errors = validate.Before(input, body_parse, opts, errors, addError)
			case "after":
				errors = validate.After(input, body_parse, opts, errors, addError)
			case "ip":
				errors = validate.Ip(input, body_parse, opts, errors, addError)
			case "password":
				errors = validate.Password(input, body_parse, opts, errors, addError)
			case "exists":
				errors = validate.Exists(input, body_parse, opts, errors, addError, dbManager)
			case "min":
				errors = validate.Min(input, body_parse, opts, errors, addError)
			case "max":
				errors = validate.Max(input, body_parse, opts, errors, addError)
			case "boolean":
				errors = validate.Boolean(input, body_parse, opts, errors, addError)
			case "sometimes":
				if _, exists_input := body_parse[input]; !exists_input {
					//Si no existe el input no se validan lás demás reglas existentes
					// log.Println("Input", input, "no existe en el body, no se validan las demás reglas")
					skipRulesMap = true
				}
			case "required":
				errors = validate.Required(input, body_parse, opts, errors, addError)
			case "required_with":
				errors = validate.RequiredWith(input, body_parse, opts, errors, addError)
			case "required_with_all":
				errors = validate.RequiredWithAll(input, body_parse, opts, errors, addError)
			case "required_without":
				errors = validate.RequiredWithout(input, body_parse, opts, errors, addError)
			case "required_without_all":
				errors = validate.RequiredWithoutAll(input, body_parse, opts, errors, addError)
			case "array":
				errors = validate.Array(input, body_parse, opts, errors, addError)
			case "type":
				errors = validate.Type(input, body_parse, opts, errors, addError)
			case "date":
				errors = validate.Date(input, body_parse, opts, errors, addError)

			default:
				errors = addError(input, input, errors, "The rule "+rule+" is not valid")
			}

			if skipRulesMap {
				break
			}
		}

		if _, ok := body_parse[input]; ok {
			safePayload[input] = body_parse[input]
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

		return gerr
	}

	if result != nil {
		jBytes, err := json.Marshal(safePayload)
		if err != nil {
			convertError := goerrors.ConvertError(err)
			return goerrors.NewGError("Error en la validación", goerrors.StatusBadRequest, &[]string{
				"Error al codificar el payload: " + err.Error(),
			}, convertError)
		}

		err = json.Unmarshal(jBytes, result)
		if err != nil {
			convertError := goerrors.ConvertError(err)
			return goerrors.NewGError("Error en la validación", goerrors.StatusBadRequest, &[]string{
				"El payload válido no se pudo convertir al tipo esperado: " + err.Error(),
			}, convertError)
		}
	}

	return nil
}

func getRuleWithOptions(rule string) []string {
	parts := strings.Split(rule, ":")

	var build_rule []string
	build_rule = append(build_rule, parts[0])

	if len(parts) > 1 {
		options := strings.Split(parts[1], ",")

		build_rule = append(build_rule, options...)
	}

	return build_rule
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
