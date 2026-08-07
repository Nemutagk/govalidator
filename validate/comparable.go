package validate

import (
	"fmt"
	"strconv"
	"time"
)

const defaultCompareDateLayout = "2006-01-02"

func toFloat64(value any) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

func parseFloat(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

type comparablePair struct {
	isDate  bool
	ownNum  float64
	cmpNum  float64
	ownDate time.Time
	cmpDate time.Time
}

// resolveComparable determina si el valor a validar es numérico o una fecha, y
// resuelve el valor contra el que se compara (otro campo del payload o un valor
// literal), soportando layouts de fecha distintos para cada lado (options[1]
// para el propio campo, options[2] para el campo/valor comparado).
func resolveComparable(input string, value any, payload map[string]any, options []string, ruleName string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) (comparablePair, map[string]interface{}, bool) {
	compareTarget := options[0]

	tmpErrorKey := fmt.Sprintf("%s.%s", input, ruleName)
	customError, hasCustomError := customeErrors[tmpErrorKey]

	fail := func(msg string) (comparablePair, map[string]interface{}, bool) {
		tmpError := msg
		if hasCustomError {
			tmpError = customError
		}
		errors = addError(input, ruleName, errors, tmpError)
		return comparablePair{}, errors, false
	}

	if numValue, ok := toFloat64(value); ok {
		var cmpValue float64
		if fieldValue, exists := payload[compareTarget]; exists {
			parsed, fieldOk := toFloat64(fieldValue)
			if !fieldOk {
				return fail(fmt.Sprintf("El campo %s no es un número válido para comparar", compareTarget))
			}
			cmpValue = parsed
		} else {
			parsed, err := parseFloat(compareTarget)
			if err != nil {
				return fail(fmt.Sprintf("El valor a comparar %s no es un número válido", compareTarget))
			}
			cmpValue = parsed
		}
		return comparablePair{isDate: false, ownNum: numValue, cmpNum: cmpValue}, errors, true
	}

	strValue, isString := value.(string)
	if !isString {
		msg := fmt.Sprintf("El campo %s debe ser un número o una fecha válida", input)
		if sliceIndex != "" {
			msg = fmt.Sprintf("El campo %s en la posición %s debe ser un número o una fecha válida", input, sliceIndex)
		}
		return fail(msg)
	}

	ownLayout := defaultCompareDateLayout
	if len(options) >= 2 && options[1] != "" {
		ownLayout = options[1]
	}

	ownDate, err := time.Parse(ownLayout, strValue)
	if err != nil {
		msg := fmt.Sprintf("El campo %s debe ser un número o una fecha válida con el formato \"%s\"", input, ownLayout)
		if sliceIndex != "" {
			msg = fmt.Sprintf("El campo %s en la posición %s debe ser un número o una fecha válida con el formato \"%s\"", input, sliceIndex, ownLayout)
		}
		return fail(msg)
	}

	targetLayout := ownLayout
	if len(options) >= 3 && options[2] != "" {
		targetLayout = options[2]
	}

	cmpDateStr := compareTarget
	if fieldValue, exists := payload[compareTarget]; exists {
		asString, ok := fieldValue.(string)
		if !ok {
			return fail(fmt.Sprintf("El campo %s no es una fecha válida para comparar", compareTarget))
		}
		cmpDateStr = asString
	}

	cmpDate, err := time.Parse(targetLayout, cmpDateStr)
	if err != nil {
		return fail(fmt.Sprintf("El valor a comparar \"%s\" no es una fecha válida con el formato \"%s\"", cmpDateStr, targetLayout))
	}

	return comparablePair{isDate: true, ownDate: ownDate, cmpDate: cmpDate}, errors, true
}
