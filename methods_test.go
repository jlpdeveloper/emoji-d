package emojid

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errStr  string
	}{
		{
			name:    "Valid emoji",
			input:   "🚀",
			wantErr: false,
		},
		{
			name:    "Invalid emoji",
			input:   "abc",
			wantErr: true,
			errStr:  "value is not in curated list of emoji",
		},
		{
			name:    "Empty string",
			input:   "",
			wantErr: true,
			errStr:  "value must not be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errStr {
				t.Errorf("New() error = %v, wantErrStr %v", err, tt.errStr)
			}
			if !tt.wantErr && got.value != tt.input {
				t.Errorf("New() got = %v, want %v", got.value, tt.input)
			}
		})
	}
}

func TestMust(t *testing.T) {
	t.Run("Valid emoji", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Must() panicked with %v", r)
			}
		}()
		id := Must("🔥")
		if id.value != "🔥" {
			t.Errorf("Must() got = %v, want %v", id.value, "🔥")
		}
	})

	t.Run("Invalid emoji panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Must() did not panic")
			}
		}()
		Must("abc")
	})
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errStr  string
	}{
		{
			name:    "Valid emoji",
			input:   "🍕",
			wantErr: false,
		},
		{
			name:    "Invalid emoji",
			input:   "xyz",
			wantErr: true,
			errStr:  "value is not in curated list of emoji",
		},
		{
			name:    "Empty string",
			input:   "",
			wantErr: true,
			errStr:  "value must not be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errStr {
				t.Errorf("Validate() error = %v, wantErrStr %v", err, tt.errStr)
			}
		})
	}
}
