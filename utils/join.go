package utils

import "strings"

func JoinURL(baseURL, endpoint string) string {
	baseURL = strings.TrimRight(baseURL, "/")
	endpoint = strings.TrimLeft(endpoint, "/")
	return baseURL + "/" + endpoint
}
