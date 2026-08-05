package emojid

import (
	"go/types"
	"testing"
)

func TestID_Scan(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantErr bool
		errStr  string
	}{
		{
			name:    "Scan string valid",
			input:   "🚀",
			wantErr: false,
		},
		{
			name:    "Scan string invalid",
			input:   "abc",
			wantErr: true,
			errStr:  "value is not in curated list of emoji",
		},
		{
			name:    "Scan byte slice valid",
			input:   []byte("🔥"),
			wantErr: false,
		},
		{
			name:    "Scan byte slice invalid",
			input:   []byte("pizza"),
			wantErr: true,
			errStr:  "value is not in curated list of emoji",
		},
		{
			name:    "Scan unsupported type",
			input:   123,
			wantErr: true,
			errStr:  "typeof value is unsupported",
		},
		{
			name:    "Scan nil",
			input:   nil,
			wantErr: true,
			errStr:  "typeof value is unsupported",
		},
		{
			name:    "Scan types.Nil",
			input:   types.Nil{},
			wantErr: true,
			errStr:  "value cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &ID{}
			err := id.Scan(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ID.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errStr {
				t.Errorf("ID.Scan() error = %v, wantErrStr %v", err, tt.errStr)
			}
			if !tt.wantErr {
				expected := ""
				switch v := tt.input.(type) {
				case string:
					expected = v
				case []byte:
					expected = string(v)
				}
				if id.value != expected {
					t.Errorf("ID.Scan() value not set correctly: got %v, want %v", id.value, expected)
				}
			}
		})
	}
}
