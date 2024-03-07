package shorten

import (
	"net/url"
	"slices"
	"strings"
)

const alphabet = "ynAJfoSgdXHB5VasEMtcbPCr1uNZ4LG723ehWkvwYR6KpxjTm8iQUFqz9D"

const lenAlphabet = uint32(len(alphabet))

func Shorten(id uint32) string {
	var digits []uint32
	var num = id

	for num > 0 {
		digits = append(digits, num%lenAlphabet)
		num /= lenAlphabet
	}

	slices.Reverse(digits)

	var builder strings.Builder
	for _, digit := range digits {
		builder.WriteString(string(alphabet[digit]))
	}

	return builder.String()
}

func PrependBaseUrl(baseUrl, id string) (string, error) {
	parsed, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	parsed.Path = id

	return parsed.String(), nil
}
