package utils

import (
	"net"
	"net/http"
	"strings"
)

func GetIP(request *http.Request) string {
	xff := request.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		ip := strings.TrimSpace(ips[0])
		if parsed := net.ParseIP(ip); parsed != nil {
			return parsed.String()
		}
	}

	xrip := request.Header.Get("X-Real-Ip")
	if xrip != "" {
		if parsed := net.ParseIP(xrip); parsed != nil {
			return parsed.String()
		}
	}

	ip, _, err := net.SplitHostPort(request.RemoteAddr)
	if err == nil {
		if parsed := net.ParseIP(ip); parsed != nil {
			return parsed.String()
		}
	}

	return ""
}
