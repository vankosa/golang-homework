package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func IsInList(myRune rune, listRunes []rune) bool {
	flag := false
	for _, val := range listRunes {
		if val == myRune {
			flag = true
			break
		}
	}
	return flag
}

func Unpack(msg string) (string, error) {
	var buf string
	var isPrevEsc bool
	resultString := ""

	for id, val := range msg {
		isPrevEsc = false
		if id == 0 {
			if unicode.IsDigit(val) {
				return "", ErrInvalidString
			}
		}
		if unicode.IsLetter(val) {
			if len(buf) != 0 {
				resultString += buf
				buf = string(val)
				continue
			}
			buf += string(val)
			continue
		}
		if unicode.IsDigit(val) {
			if isPrevEsc {
				buf += string(val)
				continue
			}
			if len(buf) == 0 {
				return "", ErrInvalidString
			}
			resultString += strings.Repeat(buf, int(val-'0'))
			buf = ""
			continue
		}
		return "", ErrInvalidString
	}

	if len(buf) != 0 {
		resultString += buf
	}
	return resultString, nil
}
