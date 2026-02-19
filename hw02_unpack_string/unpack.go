package hw02unpackstring

import (
	"errors"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	strRunes := []rune(str)
	var currentRune rune
	for _, c := range strRunes {
		if currentRune == 0 {
			if unicode.IsDigit(c) {
				return "", ErrInvalidString
			} else {
				currentRune = c
			}
		}

	}
	return "", nil
}
