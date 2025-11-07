package govalidator

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Nemutagk/godb"
	"github.com/Nemutagk/godb/definitions/db"
	"github.com/Nemutagk/goerrors"
	"github.com/Nemutagk/govalidator/validate"
)

type validInput struct {
	Name  string
	Rules []validRule
}

type validRule struct {
	Rule string
	Opts *[]string
}

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

	all_rules := []validInput{}

	for input, rule := range rules {
		separated_rules := strings.Split(rule, "|")

		struct_input := validInput{
			Name:  input,
			Rules: []validRule{},
		}

		for _, rule := range separated_rules {
			tmp_rule := getRuleWithOptions(rule)
			opts := tmp_rule[1:]
			struct_input.Rules = append(struct_input.Rules, validRule{
				Rule: tmp_rule[0],
				Opts: &opts,
			})
		}
		all_rules = append(all_rules, struct_input)
	}

	errors := make(map[string]any)
	safePayload := make(map[string]any)

	dbManager := godb.InitConnections(map[string]db.DbConnection{})

	for _, input := range all_rules {

		skipRulesMap := false

		for _, rule := range input.Rules {
			opts := []string{}

			if rule.Opts != nil {
				opts = *rule.Opts
			}

			switch rule.Rule {
			case "email":
				errors = validate.Email(input.Name, body_parse, opts, errors, addError)
			case "confirmation":
				errors = validate.Confirmation(input.Name, body_parse, opts, errors, addError)
			case "unique":
				errors = validate.Unique(input.Name, body_parse, opts, errors, addError, dbManager)
			case "in":
				errors = validate.In(input.Name, body_parse, opts, errors, addError)
			case "before":
				errors = validate.Before(input.Name, body_parse, opts, errors, addError)
			case "after":
				errors = validate.After(input.Name, body_parse, opts, errors, addError)
			case "ip":
				errors = validate.Ip(input.Name, body_parse, opts, errors, addError)
			case "password":
				errors = validate.Password(input.Name, body_parse, opts, errors, addError)
			case "exists":
				errors = validate.Exists(input.Name, body_parse, opts, errors, addError, dbManager)
			case "min":
				errors = validate.Min(input.Name, body_parse, opts, errors, addError)
			case "max":
				errors = validate.Max(input.Name, body_parse, opts, errors, addError)
			case "boolean":
				errors = validate.Boolean(input.Name, body_parse, opts, errors, addError)
			case "sometimes":
				if _, exists_input := body_parse[input.Name]; !exists_input {
					//Si no existe el input no se validan lás demás reglas existentes
					// log.Println("Input", input, "no existe en el body, no se validan las demás reglas")
					skipRulesMap = true
				}
			case "required":
				errors = validate.Required(input.Name, body_parse, opts, errors, addError)
			case "required_with":
				errors = validate.RequiredWith(input.Name, body_parse, opts, errors, addError)
			case "required_with_all":
				errors = validate.RequiredWithAll(input.Name, body_parse, opts, errors, addError)
			case "required_without":
				errors = validate.RequiredWithout(input.Name, body_parse, opts, errors, addError)
			case "required_without_all":
				errors = validate.RequiredWithoutAll(input.Name, body_parse, opts, errors, addError)
			case "array":
				errors = validate.Array(input.Name, body_parse, opts, errors, addError)
			case "type":
				errors = validate.Type(input.Name, body_parse, opts, errors, addError)
			case "date":
				errors = validate.Date(input.Name, body_parse, opts, errors, addError)
			case "nullable":
				isNull := false
				errors, isNull = validate.Nullable(input.Name, body_parse, opts, errors, addError)

				if isNull {
					skipRulesMap = true
				}
			default:
				errors = addError(input.Name, input.Name, errors, fmt.Sprintf("The rule %s is not valid", rule.Rule))
			}

			if skipRulesMap {
				break
			}
		}

		if _, ok := body_parse[input.Name]; ok {
			safePayload[input.Name] = body_parse[input.Name]
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
	if !strings.Contains(rule, ":") {
		return []string{rule}
	}

	parts := strings.SplitN(rule, ":", 2)

	opts := strings.Split(parts[1], ",")

	all := []string{parts[0]}
	all = append(all, opts...)

	return all
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
