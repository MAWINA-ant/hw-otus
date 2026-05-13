package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	NamedResponse struct {
		Name string
		Resp Response `validate:"nested"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: Token{
				Header:    []byte("fdsafd"),
				Payload:   []byte("fdscxzvvcxzafd"),
				Signature: []byte("34243243234"),
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "qwertyuiopasdfghjklzxcvbnm1234567890",
				Name:   "Ruslan",
				Age:    12,
				Email:  "cfadsfasd@fdadsa.com",
				Role:   "developer",
				Phones: []string{"89159532200", "845612365462"},
				meta:   json.RawMessage{},
			},
			expectedErr: ValidationErrors{
				CreateValidationError("Age", "min"),
				CreateValidationError("Role", "in"),
				CreateValidationError("Phones", "len"),
			},
		},
		{
			in: Response{
				Code: 400,
			},
			expectedErr: ValidationErrors{
				CreateValidationError("Code", "in"),
			},
		},
		{
			in: User{
				ID:     "qwertyuiopasdfghjklzxcvbnm1234567890",
				Name:   "Ruslan",
				Age:    25,
				Email:  "cfadsfasd@fdadsa.com",
				Role:   "stuff",
				Phones: []string{"89159532200", "84561236546"},
				meta:   json.RawMessage{},
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1234567",
			},
			expectedErr: ValidationErrors{
				CreateValidationError("Version", "len"),
			},
		},
		{
			in: NamedResponse{
				Name: "BlaBlaBla",
				Resp: Response{
					Code: 401,
					Body: "foobar",
				},
			},
			expectedErr: ValidationErrors{
				CreateValidationError("Code", "in"),
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
			_ = tt
		})
	}
}
