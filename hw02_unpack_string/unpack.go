package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}

	sb := strings.Builder{}
	runes := []rune(str)

	for i, r := range runes {
		// Является ли символ последним в строке.
		if i == len(runes)-1 {
			if !unicode.IsDigit(r) {
				sb.WriteRune(r)
				break
			}
		}
		// Является ли символ цифрой.
		if unicode.IsDigit(r) {
			// Является ли первый или следующий символ цифрой.
			if i == 0 || unicode.IsDigit(runes[i+1]) {
				return "", ErrInvalidString
			}
		} else {
			if unicode.IsDigit(runes[i+1]) {
				// Преобразуем символ в цифру.
				count, _ := strconv.Atoi(string(runes[i+1]))
				for j := 0; j < count; j++ {
					sb.WriteRune(r)
				}
				continue
			}
			sb.WriteRune(r)
		}
	}
	return sb.String(), nil
}
