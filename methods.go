package emojid

// New creates a new ID from a string value. It returns an error if the value
// is not a valid emoji ID.
func New(value string) (ID, error) {
	err := validateStr(value)
	if err != nil {
		return ID{}, err
	}
	return ID{
		value: value,
	}, nil
}

// Must creates a new ID from a string value. It panics if the value is not
// a valid emoji ID.
func Must(value string) ID {
	i, err := New(value)
	if err != nil {
		panic(err)
	}
	return i
}

// Validate checks if the given string is a valid emoji ID.
func Validate(value string) error {
	return validateStr(value)
}
