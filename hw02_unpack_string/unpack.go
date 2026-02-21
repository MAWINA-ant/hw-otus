package hw02unpackstring

import (
	"errors"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var resultRunes []rune
	canMultiply := false
	isBackSlash := false
	for _, v := range str {
		if unicode.IsDigit(v) {
			if isBackSlash {
				resultRunes = append(resultRunes, v)
				canMultiply = true
				isBackSlash = false
			} else {
				if !canMultiply {
					return "", ErrInvalidString
				}
				countRepetitions := v - '0'
				if countRepetitions == 0 {
					resultRunes = resultRunes[:len(resultRunes)-1]
				} else {
					appendRune := resultRunes[len(resultRunes)-1]
					for i := 0; i < int(countRepetitions)-1; i++ {
						resultRunes = append(resultRunes, appendRune)
					}
				}
				canMultiply = false
			}
		} else if v == '\\' {
			if isBackSlash {
				resultRunes = append(resultRunes, v)
				canMultiply = true
				isBackSlash = false
			} else {
				isBackSlash = true
			}
		} else {
			if isBackSlash {
				return "", ErrInvalidString
			}
			resultRunes = append(resultRunes, v)
			canMultiply = true
		}
	}
	return string(resultRunes), nil
}
