package scraper

import (
	"regexp"
	"strings"
)

var escapeSeqRegex = regexp.MustCompile(`[\n\r\t]+`)

func cleanEscapedText(value string) string {
	return escapeSeqRegex.ReplaceAllString(value, "")
}

func parseLocationName(value string) string {
	runes := []rune(value)
	if len(runes) < 3 {
		return value + "*"
	}
	short := string(runes[:3])
	if strings.Contains(short, "*") {
		return short
	}
	return short + "*"
}
