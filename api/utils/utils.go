package utils

import (
	"os"
	"strings"

	neturl "net/url"
)

func IsDifferentDomain(u string) bool {
	parsedURL, err := neturl.Parse(u)
	if err != nil {
		return false
	}
	return parsedURL.Host != os.Getenv("DOMAIN")
}

func EnsureHttpPrefix(u string) string {
	if !strings.HasPrefix(u, "http") || !strings.HasPrefix(u, "https") {
		return "http://" + u
	}
	return u
}
