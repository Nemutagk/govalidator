package validate

import (
	"fmt"
	"net"
	"strings"
)

func Ip(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if value == nil || value == "" {
		return errors
	}

	ip := value.(string)
	if !strings.Contains(ip, ",") {
		if net.ParseIP(strings.TrimSpace(ip)) == nil {
			tmpError := "La dirección IP " + ip + " no es una dirección IP válida"

			if sliceIndex != "" {
				tmpError = fmt.Sprintf("La dirección IP %s en la posición %s no es válida", ip, sliceIndex)
			}

			customeErrorKey := fmt.Sprintf("%s.ip", input)
			if customeError, exists := customeErrors[customeErrorKey]; exists {
				tmpError = customeError
			}
			errors = addError(input, "ip", errors, tmpError)
		}
	} else {
		ipList := strings.Split(ip, ",")
		for index, ip := range ipList {
			if net.ParseIP(strings.TrimSpace(ip)) == nil {
				tmpError := fmt.Sprintf("La dirección IP %d:%s no es válida", index, ip)

				if sliceIndex != "" {
					tmpError = fmt.Sprintf("La dirección IP %s en la posición %s no es válida", ip, sliceIndex)
				}

				customeErrorKey := fmt.Sprintf("%s.ip", input)
				if customeError, exists := customeErrors[customeErrorKey]; exists {
					tmpError = customeError
				}
				errors = addError(input, "ip", errors, tmpError)
			}
		}
	}

	return errors
}
