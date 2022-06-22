package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	sr := []rune(str)

	var sb strings.Builder

	for i, char := range sr {
		if unicode.IsDigit(char) && i == 0 {
			return "", ErrInvalidString
		}

		if unicode.IsDigit(char) && unicode.IsDigit(sr[i-1]) {
			return "", ErrInvalidString
		}

		if unicode.IsDigit(char) {
			n, _ := strconv.Atoi(string(char))
			if n == 0 {
				tr := []rune(sb.String())
				tr = tr[:len(tr)-1]

				sb.Reset()
				sb.WriteString(string(tr))

				continue
			}

			repeater := strings.Repeat(string(sr[i-1]), n-1)
			sb.WriteString(repeater)
			continue
		}

		sb.WriteRune(char)
	}

	return sb.String(), nil
}
