package emojid

import (
	"errors"
	"go/types"
)

type ID struct {
	value string
}

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
