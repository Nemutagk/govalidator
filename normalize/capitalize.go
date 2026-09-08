package normalize

import "unicode"

func Capitalize(value any, options []string) any {
	s, ok := value.(string)
	if !ok || s == "" {
		return value
	}

	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])

	return string(r)
}
