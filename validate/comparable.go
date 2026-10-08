package validate

import (
	"cmp"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

const defaultCompareDateLayout = "2006-01-02"

// toFloat64 convierte cualquier tipo numérico de Go (int*, uint*, float*, incluidos los
// tipos con nombre) a float64. Un float32 se convierte por su representación decimal más
// corta para que 0.1 siga siendo 0.1 y no 0.10000000149011612.
func toFloat64(value any) (float64, bool) {
	rv := reflect.ValueOf(value)

	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(rv.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(rv.Uint()), true
	case reflect.Float32:
		parsed, err := strconv.ParseFloat(strconv.FormatFloat(rv.Float(), 'g', -1, 32), 64)
		return parsed, err == nil
	case reflect.Float64:
		return rv.Float(), true
	default:
		return 0, false
	}
}

// scalarToString convierte un valor escalar (string, bool, entero, uint o flotante, incluidos
// los tipos con nombre) a su texto decimal. ok es false para nil, slices, maps y structs.
func scalarToString(value any) (text string, ok bool) {
	rv := reflect.ValueOf(value)

	switch rv.Kind() {
	case reflect.String:
		return rv.String(), true
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool()), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10), true
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 32), true
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 64), true
	default:
		return "", false
	}
}

// equalsOption compara un valor contra una opción de regla (siempre un string) como texto
// contra texto; nil, slices, maps y structs nunca son iguales.
func equalsOption(value any, option string) bool {
	text, ok := scalarToString(value)
	return ok && text == option
}

// compareToInt64 compara cualquier tipo numérico de Go contra un entero y devuelve -1, 0
// o 1 (valor menor, igual o mayor que limit). ok es false si el valor no es numérico.
// Los enteros se comparan como int64 sin pasar por float64, y un uint64 que no cabe en
// int64 siempre es mayor que limit.
func compareToInt64(value any, limit int64) (result int, ok bool) {
	rv := reflect.ValueOf(value)

	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return cmp.Compare(rv.Int(), limit), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if rv.Uint() > math.MaxInt64 {
			return 1, true
		}
		return cmp.Compare(int64(rv.Uint()), limit), true
	case reflect.Float32, reflect.Float64:
		floatValue, _ := toFloat64(value)
		return cmp.Compare(floatValue, float64(limit)), true
	default:
		return 0, false
	}
}

func parseFloat(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

// resolveNumericTarget resuelve el valor numérico contra el que se compara: otro campo del
// payload si existe una clave con ese nombre, o el literal numérico en caso contrario.
// errMsg viene vacío cuando se pudo resolver.
func resolveNumericTarget(payload map[string]any, target string) (value float64, isField bool, errMsg string) {
	if fieldValue, exists := payload[target]; exists {
		parsed, ok := toFloat64(fieldValue)
		if !ok {
			return 0, true, fmt.Sprintf("El campo %s no es un número válido para comparar", target)
		}
		return parsed, true, ""
	}

	parsed, err := parseFloat(target)
	if err != nil {
		return 0, false, fmt.Sprintf("El valor a comparar %s no es un número válido", target)
	}
	return parsed, false, ""
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
