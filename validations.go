package emojid

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// validateStr checks if the given string is a valid emoji ID by ensuring
// it is not empty, is valid UTF-8, and contains only runes from the
// curated emoji list.
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
