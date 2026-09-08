package normalize

import "strings"

func Trim(value any, options []string) any {
	s, ok := value.(string)
	if !ok {
		return value
	}

	return strings.TrimSpace(s)
}
