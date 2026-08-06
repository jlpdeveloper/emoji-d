package emojid

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"go/types"
)

// ID represents a unique identifier validated against a curated list of emojis.
type ID struct {
	value string
}

// Scan implements the sql.Scanner interface. It validates that the value
// is a string or []byte and matches the curated emoji list.
func (i *ID) Scan(value any) error {

	switch v := value.(type) {
	case string:
		if err := validateStr(v); err != nil {
			return err
		}
		i.value = v
	case []byte:
		s := string(v)
		if err := validateStr(s); err != nil {
			return err
		}
		i.value = s
	case types.Nil:
		return errors.New("value cannot be nil")
	default:
		return errors.New("typeof value is unsupported")
	}
	return nil
}

// Value implements the driver.Valuer interface. It returns the string value
// or an error if the value is not defined.
func (i *ID) Value() (driver.Value, error) {
	if i.value == "" {
		return nil, errors.New("value is not defined")
	}
	return i.value, nil
}

// String returns the string representation of the ID.
func (i *ID) String() string {
	return i.value
}

// MarshalJSON implements the json.Marshaler interface. It returns the
// JSON encoding of the ID.
func (i *ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.value)
}

// UnmarshalJSON implements the json.Unmarshaler interface. It validates
// the JSON-encoded data and stores the result in the ID.
func (i *ID) UnmarshalJSON(data []byte) error {
	v := new("")
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	err = validateStr(*v)
	if err != nil {
		return err
	}
	i.value = *v
	return nil
}

// IsValid returns true if the ID has a non-empty value.
func (i *ID) IsValid() bool {
	return i.value != ""
}
