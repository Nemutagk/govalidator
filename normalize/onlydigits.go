package normalize

import "strings"

func OnlyDigits(value any, options []string) any {
	s, ok := value.(string)
	if !ok {
		return value
	}

	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}

	return b.String()
}
