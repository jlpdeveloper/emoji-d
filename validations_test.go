package emojid

import (
	"testing"
)

func Test_validateStr(t *testing.T) {
	tests := []struct {
		wantErr bool
		input   string
		errStr  string
		name    string
	}{
		{
			name:    "Successful rune parsing",
			wantErr: false,
			input:   "🦆",
		},
		{
			name:    "Empty string",
			input:   "",
			wantErr: true,
			errStr:  "value must not be empty",
		},
		{
			name:    "Only spaces",
			input:   "        ",
			wantErr: true,
			errStr:  "value must not be empty",
		},
		{
			name:    "invalid utf8 string",
			input:   string([]byte{0xff}),
			wantErr: true,
			errStr:  "value must be utf8",
		},
		{
			name:    "Multiple valid emojis",
			input:   "🚀🔥",
			wantErr: false,
		},
		{
			name:    "Emoji not in curated list",
			input:   "🌍",
			wantErr: true,
			errStr:  "value is not in curated list of emoji",
		},
		{
			name:    "ASCII characters",
			input:   "abc",
			wantErr: true,
			errStr:  "value is not in curated list of emoji",
		},
		{
			name:    "Invalid rune (surrogate half)",
			input:   string([]byte{0xED, 0xA0, 0x80}),
			wantErr: true,
			errStr:  "value must be utf8",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateStr(test.input)
			if (err != nil) != test.wantErr {
				t.Errorf("wanted error %v, got err %v", test.wantErr, err)
			}
			if test.wantErr {
				if err != nil && err.Error() != test.errStr {
					t.Errorf("wanted error %v, got error %v", test.errStr, err.Error())
				}
			}
		})
	}
}
