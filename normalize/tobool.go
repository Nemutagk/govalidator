package normalize

import "strings"

func ToBool(value any, options []string) any {
	s, ok := value.(string)
	if !ok {
		return value
	}

	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		return value
	}
}
