package normalize

import (
	"strconv"
	"strings"
)

func ToFloat(value any, options []string) any {
	switch v := value.(type) {
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return value
		}
		return f
	case int:
		return float64(v)
	case float64:
		return v
	default:
		return value
	}
}
