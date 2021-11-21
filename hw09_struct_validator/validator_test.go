package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"fmt"
	"testing"
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
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name:        "string input",
			in:          "test string",
			expectedErr: ErrInvalidArg,
		},
		{
			name: "valid user",
			in: User{
				ID:     "MznkVp0AIIpPrZyyzoYXgj5M3XY1pjbbp5kb",
				Name:   "Ivan",
				Age:    18,
				Email:  "test@gmail.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: nil,
		},
		{
			name: "user with invalid ID",
			in: User{
				ID:     "1234554345",
				Name:   "Vasiliy",
				Age:    18,
				Email:  "test@gmail.com",
				Role:   "stuff",
				Phones: []string{"12345678901"},
			},
			expectedErr: &ValidationErrors{
				{
					Field: "ID",
					Err:   ErrInvalidLength,
				},
			},
		},
		{
			name: "user with invalid phone",
			in: User{
				ID:     "aznkVp0AIIpPrZyyzoYXgj5M3XY1pjbbp5ka",
				Name:   "Andrey",
				Age:    18,
				Email:  "test@gmail.com",
				Role:   "admin",
				Phones: []string{"123"},
			},
			expectedErr: &ValidationErrors{
				{
					Field: "Phones",
					Err:   ErrInvalidLength,
				},
			},
		},
		{
			name: "user with invalid ID small age invalid email role phones",
			in: User{
				ID:     "1234554345",
				Name:   "Andrey",
				Age:    17,
				Email:  "invalid-email",
				Role:   "worker",
				Phones: []string{"1234567"},
			},
			expectedErr: &ValidationErrors{
				{
					Field: "ID",
					Err:   ErrInvalidLength,
				},
				{
					Field: "Age",
					Err:   ErrValueLessMin,
				},
				{
					Field: "Email",
					Err:   ErrPatternNotMatch,
				},
				{
					Field: "Role",
					Err:   ErrValueNotIn,
				},
				{
					Field: "Phones",
					Err:   ErrInvalidLength,
				},
			},
		},
		{
			name: "user with invalid ID big age phones",
			in: User{
				ID:     "1234554345",
				Name:   "Andrey",
				Age:    51,
				Email:  "test@gmail.com",
				Role:   "admin",
				Phones: []string{"1234567"},
			},
			expectedErr: &ValidationErrors{
				{
					Field: "ID",
					Err:   ErrInvalidLength,
				},
				{
					Field: "Age",
					Err:   ErrValueGreaterMax,
				},
				{
					Field: "Phones",
					Err:   ErrInvalidLength,
				},
			},
		},
		{
			name: "valid app",
			in: App{
				Version: "12345",
			},
			expectedErr: nil,
		},
		{
			name: "app with invalid version",
			in: App{
				Version: "1234",
			},
			expectedErr: &ValidationErrors{
				{
					Field: "Version",
					Err:   ErrInvalidLength,
				},
			},
		},
		{
			name: "token",
			in: Token{
				Header:  []byte("test-header"),
				Payload: []byte("test-payload"),
			},
			expectedErr: nil,
		},
		{
			name: "valid response",
			in: Response{
				Code: 404,
				Body: "",
			},
			expectedErr: nil,
		},
		{
			name: "response with invalid code",
			in: Response{
				Code: 123,
				Body: "",
			},
			expectedErr: &ValidationErrors{
				{
					Field: "Code",
					Err:   ErrValueNotIn,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.name), func(t *testing.T) {
			actual := Validate(tt.in)
			if tt.expectedErr != nil && actual == nil ||
				tt.expectedErr == nil && actual != nil {
				t.Errorf(
					"got: %v, expected: %s",
					actual, tt.expectedErr,
				)
			}

			if tt.expectedErr != nil && actual != nil && tt.expectedErr.Error() != actual.Error() {
				t.Errorf(
					"got: %v, expected: %s",
					actual, tt.expectedErr,
				)
			}
		})
	}
}
