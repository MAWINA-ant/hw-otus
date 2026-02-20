package hw02unpackstring

import (
	"errors"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var resultRunes []rune
	canAppend := false
	for _, c := range str {
		if unicode.IsDigit(c) {
			if !canAppend {
				return "", ErrInvalidString
			}
			countRepetitions := c - '0'
			if countRepetitions == 0 {
				resultRunes = resultRunes[:len(resultRunes)-1]
			} else {
				appendRune := resultRunes[len(resultRunes)-1]
				for i := 0; i < int(countRepetitions)-1; i++ {
					resultRunes = append(resultRunes, appendRune)
				}
			}
			canAppend = false
		} else {
			resultRunes = append(resultRunes, c)
			canAppend = true
		}
	}
	return string(resultRunes), nil
}
