package main

import (
	"regexp"
	"strings"
)

func trimSpacesAndTabs(str string) string {
	r := regexp.MustCompile("\\s+")
	replace := r.ReplaceAllString(str, " ")
	return strings.TrimSpace(replace)
}

// удаление возможного символа \r в конце строки (c'est la Windows)
func removeR(str string) string {
	if len(str) > 1 {
		if str[len(str)-1] == '\r' {
			str = str[:len(str)-1]
		}
	}
	return str
}
