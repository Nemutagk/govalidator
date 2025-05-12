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
			errors = addError(input, "ip", errors, "The input "+input+" is not a valid IP address 1")
		}
	} else {
		ipList := strings.Split(ip, ",")
		for _, ip := range ipList {
			if net.ParseIP(strings.TrimSpace(ip)) == nil {
				errors = addError(input, "ip", errors, "The ip "+ip+" is not a valid IP address 2")
			}
		}
	}

	return errors
}
