package emojid

import (
	"errors"
	"strings"
	"unicode/utf8"
)

func validateStr(s string) error {
	if strings.TrimSpace(s) == "" {
		return errors.New("value must not be empty")
	}
	if !utf8.ValidString(s) {
		return errors.New("value must be utf8")
	}
	runes := []rune(s)
	for _, r := range runes {
		if !utf8.ValidRune(r) {
			return errors.New("value contains invalid runes")
		}
		if _, ok := validRunes[r]; !ok {
			return errors.New("value is not in curated list of emoji")
		}
	}
	return nil
}
