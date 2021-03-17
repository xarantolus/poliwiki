package util

import (
	"strings"
	"unicode"
)

// https://help.twitter.com/de/using-twitter/how-to-use-hashtags
func Hashtag(lastname string) string {
	f := strings.FieldsFunc(lastname, func(r rune) bool {
		return !unicode.IsLetter(r)
	})

	if len(f) > 1 {
		lastname = strings.Join(f, "")
	}

	return "#" + lastname
}
