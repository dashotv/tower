package app

import (
	"regexp"
	"strings"
)

var pathQuoteRegex = regexp.MustCompile(`'(\w{1,2})`)
var pathCharRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func path(title string) string {
	var s string
	s = pathQuoteRegex.ReplaceAllString(title, "$1")
	s = strings.ToLower(s)
	s = pathCharRegex.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	return s
}
