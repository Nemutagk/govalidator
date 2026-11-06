package helper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

func GenerateUuid() string {
	buil, err := uuid.NewV7()
	if err != nil {
		return ""
	}

	return buil.String()
}

func PrettyPrint(data any) {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}

	fmt.Println(string(prettyJSON))
}

func SliceContains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func IsEmpty(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)

	// Desenvuelve punteros e interfaces
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return true
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.String:
		return strings.TrimSpace(rv.String()) == ""
	case reflect.Slice, reflect.Array, reflect.Map:
		return rv.Len() == 0
	case reflect.Bool:
		return false // normalmente false NO se considera "vacío"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return false // si quieres, aquí podrías considerar 0 como vacío
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return false
	case reflect.Float32, reflect.Float64:
		return false
	case reflect.Struct:
		// opcional: comparar con zero-value del struct
		zero := reflect.Zero(rv.Type())
		return reflect.DeepEqual(rv.Interface(), zero.Interface())
	default:
		return false
	}
}

// AsStringMap devuelve un map con claves string o enteras (map[string]string, map[int]int...)
// como map[string]any, con las claves en texto, para recorrerlo con notación de punto. Un
// map[string]any se devuelve tal cual (sin copiar) para que las escrituras sigan llegando al
// original. ok es false si el valor no es un map o su clave es de otro tipo.
// shortcut: una clave que contiene un punto no se puede direccionar con 'a.b'.
func AsStringMap(value any) (map[string]any, bool) {
	if m, ok := value.(map[string]any); ok {
		return m, true
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Map {
		return nil, false
	}
	if k := rv.Type().Key().Kind(); k != reflect.String && (k < reflect.Int || k > reflect.Uint64) {
		return nil, false
	}

	result := make(map[string]any, rv.Len())
	for iter := rv.MapRange(); iter.Next(); {
		result[fmt.Sprint(iter.Key().Interface())] = iter.Value().Interface()
	}
	return result, true
}
