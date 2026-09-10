package govalidator

import (
	"fmt"
	"reflect"
	"strings"
)

type StructOptions struct {
	Tag string
}

func ValidateStruct[T any](s T, inputs []Input, customeallErrors map[string]string, models map[string]func(data any, payload map[string]any, opts *[]string) (bool, string), opts ...StructOptions) (T, error) {
	var zero T

	tag := "json"
	if len(opts) > 0 && opts[0].Tag != "" {
		tag = opts[0].Tag
	}

	body, err := structToMap(s, tag)
	if err != nil {
		return zero, fmt.Errorf("error convirtiendo struct a map: %w", err)
	}

	// log.Printf("body: %+v", body)
	safePayload, err := ValidateRequest(body, inputs, customeallErrors, models)
	if err != nil {
		return zero, err
	}

	result, err := ConvertPayload[T](safePayload)
	if err != nil {
		return zero, fmt.Errorf("error convirtiendo payload a struct: %w", err)
	}

	return result, nil
}

func structToMap(s any, tag string) (map[string]any, error) {
	return reflectStructToMap(reflect.ValueOf(s), tag)
}

func reflectStructToMap(v reflect.Value, tag string) (map[string]any, error) {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, nil
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("se esperaba un struct, se recibió %s", v.Kind())
	}

	result := make(map[string]any)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		if !field.IsExported() {
			continue
		}

		key, omitempty := fieldNameAndOpts(field, tag)
		if key == "-" {
			continue
		}
		if omitempty && isEmptyValue(fieldVal) {
			continue
		}

		val, err := reflectToAny(fieldVal, tag)
		if err != nil {
			return nil, err
		}
		result[key] = val
	}

	return result, nil
}

func reflectToAny(v reflect.Value, tag string) (any, error) {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		if v.Type().String() == "time.Time" {
			return v.Interface(), nil
		}
		return reflectStructToMap(v, tag)
	case reflect.Slice:
		if v.IsNil() {
			return nil, nil
		}
		// Preservar tipo de slices primitivos (no recursivos)
		elemKind := v.Type().Elem().Kind()
		if elemKind == reflect.Ptr {
			elemKind = v.Type().Elem().Elem().Kind()
		}
		if elemKind != reflect.Struct && elemKind != reflect.Interface && elemKind != reflect.Map {
			return v.Interface(), nil
		}

		result := make([]any, v.Len())
		for i := range v.Len() {
			val, err := reflectToAny(v.Index(i), tag)
			if err != nil {
				return nil, err
			}
			result[i] = val
		}
		return result, nil
	case reflect.Map:
		if v.IsNil() {
			// fallback para mapas nulos
		}
		// Preservar mapas primitivos si es necesario, o al menos convertir controladamente
		elemKind := v.Type().Elem().Kind()
		if elemKind == reflect.Ptr {
			elemKind = v.Type().Elem().Elem().Kind()
		}
		if elemKind != reflect.Struct && elemKind != reflect.Interface && elemKind != reflect.Map {
			return v.Interface(), nil
		}

		result := make(map[string]any)
		for _, k := range v.MapKeys() {
			val, err := reflectToAny(v.MapIndex(k), tag)
			if err != nil {
				return nil, err
			}
			result[fmt.Sprintf("%v", k.Interface())] = val
		}
		return result, nil
	default:
		// Los valores con un tipo con nombre distinto de su tipo subyacente
		// (ej. "type StatusType string") conservan ese tipo con nombre al
		// pasar por v.Interface(), y una comparación por igualdad entre una
		// interfaz que envuelve StatusType("active") y una que envuelve
		// string("active") da false en Go, aunque el contenido sea idéntico
		// (la comparación de interfaces exige que el tipo dinámico también
		// coincida). Reglas como "in"/"not_in"/"equal" comparan el valor
		// contra literales de Go planos (string, bool, etc.), así que aquí
		// se "desenvuelve" el valor a su tipo primitivo subyacente para que
		// esas comparaciones funcionen igual que con un campo sin tipo con
		// nombre. Los tipos ya primitivos (string, int, bool, etc.) no se
		// tocan: v.Type().Name() coincide con v.Kind().String() para ellos.
		if v.Type().Name() != v.Kind().String() {
			switch v.Kind() {
			case reflect.String:
				return v.String(), nil
			case reflect.Bool:
				return v.Bool(), nil
			case reflect.Int:
				return int(v.Int()), nil
			case reflect.Int8:
				return int8(v.Int()), nil
			case reflect.Int16:
				return int16(v.Int()), nil
			case reflect.Int32:
				return int32(v.Int()), nil
			case reflect.Int64:
				return v.Int(), nil
			case reflect.Uint:
				return uint(v.Uint()), nil
			case reflect.Uint8:
				return uint8(v.Uint()), nil
			case reflect.Uint16:
				return uint16(v.Uint()), nil
			case reflect.Uint32:
				return uint32(v.Uint()), nil
			case reflect.Uint64:
				return v.Uint(), nil
			case reflect.Float32:
				return float32(v.Float()), nil
			case reflect.Float64:
				return v.Float(), nil
			}
		}
		return v.Interface(), nil
	}
}

func fieldNameAndOpts(field reflect.StructField, tag string) (name string, omitempty bool) {
	tagVal := field.Tag.Get(tag)
	if tagVal == "" {
		return field.Name, false
	}
	parts := strings.Split(tagVal, ",")
	name = parts[0]
	if name == "" {
		name = field.Name
	}
	for _, opt := range parts[1:] {
		if opt == "omitempty" {
			omitempty = true
		}
	}
	return name, omitempty
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	default:
		return v.IsZero()
	}
}
