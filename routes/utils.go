package routes

import "strings"

func joinURL(parts ...string) string {
	for i, part := range parts {
		parts[i] = strings.Trim(part, "/")
	}
	return strings.Join(parts, "/")
}

func replaceParam(url string, paramName string, paramValue string) string {
	return strings.ReplaceAll(url, ":"+paramName, paramValue)
}

func replaceParams(url string, params map[string]string) string {
	for paramName, paramValue := range params {
		url = replaceParam(url, paramName, paramValue)
	}
	return url
}
