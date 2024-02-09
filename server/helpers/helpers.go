package helpers

import (
	"errors"
	"strings"
	"unicode"
)


func FormatIsbn(isbn string) (string, error) {
	str := strings.Map(FilterOnlyDigits, isbn)
	strLen := len(str)
	if strLen != 10 && strLen != 13 {
		return "", errors.New("invalid isbn")
	}
	return str, nil
}

func FilterOnlyDigits(r rune) rune {
	if unicode.IsDigit(r) {
		return r
	} else {
		return -1
	}
}
