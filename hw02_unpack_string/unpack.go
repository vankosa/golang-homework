package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var (
		result                    strings.Builder
		previousRune, currentRune rune
	)

	for _, currentRune = range s {
		if previousRune == 0 {
			previousRune = currentRune
			continue
		}

		if unicode.IsDigit(currentRune) {
			if unicode.IsDigit(previousRune) {
				return "", errors.New("invalid string")
			}

			count, err := strconv.Atoi(string(currentRune))
			if err != nil {
				return "", err
			}

			_, err = result.WriteString(strings.Repeat(string(previousRune), count))
			if err != nil {
				return "", err
			}
		} else if !unicode.IsDigit(previousRune) {
			result.WriteString(string(previousRune))
		}

		previousRune = currentRune
	}

	if currentRune != 0 && !unicode.IsDigit(currentRune) {
		result.WriteString(string(currentRune))
	}

	return result.String(), nil
}
