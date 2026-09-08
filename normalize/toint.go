package normalize

import (
	"strconv"
	"strings"
)

func ToInt(value any, options []string) any {
	switch v := value.(type) {
	case string:
		i, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return value
		}
		return i
	case float64:
		return int(v)
	case int:
		return v
	default:
		return value
	}
}
