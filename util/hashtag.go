package util

import (
	"strings"
	"unicode"
)

// Hashtag returns the given text with a hashtag # in front of it.
// If there are any characters that are not valid, they will be removed,
// e.g. a name of the form Name-Name will become "#NameName"
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
