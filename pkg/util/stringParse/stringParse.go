package stringParse

import (
	"regexp"
	"strings"
)

func Reform(s string) string {
	return strings.ReplaceAll(s, "/", "")
}

func IsNsqSysMessage(s string) bool {
	re := regexp.MustCompile(`^\d{4}/\d{2}/\d{2}`)
	return re.MatchString(s)
}
