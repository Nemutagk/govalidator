package validate

import (
	"net"
	"strings"
)

func Ip(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	value := payload[input]

	if value == nil || value == "" {
		return errors
	}

	ip := value.(string)
	if !strings.Contains(ip, ",") {
		if net.ParseIP(strings.TrimSpace(ip)) == nil {
			errors = addError(input, "ip", errors, "La dirección IP "+ip+" no es una dirección IP válida")
		}
	} else {
		ipList := strings.Split(ip, ",")
		for _, ip := range ipList {
			if net.ParseIP(strings.TrimSpace(ip)) == nil {
				errors = addError(input, "ip", errors, "La dirección IP "+ip+" no es válida")
			}
		}
	}

	return errors
}
