package normalize

import "strings"

func Upper(value any, options []string) any {
	s, ok := value.(string)
	if !ok {
		return value
	}

	return strings.ToUpper(s)
}
