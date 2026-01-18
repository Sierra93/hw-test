package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
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
		in          interface{}
		expectedErr error
	}{
		{
			in:          "not a struct",
			expectedErr: ErrInvalidInputType,
		},
		{
			in: struct {
				Field string `validate:"len:tt"`
			}{},
			expectedErr: ErrInvalidInteger,
		},
		{
			in: struct {
				Field []int `validate:"min:11|max:tt"`
			}{},
			expectedErr: ErrInvalidInteger,
		},
		{
			in: struct {
				Field []int `validate:"in:11,tt"`
			}{},
			expectedErr: ErrInvalidInteger,
		},
		{
			in: struct {
				Field []int `validate:"tt:11"`
			}{},
			expectedErr: ErrInvalidParam,
		},
		{
			in: struct {
				Field bool `validate:"max:11"`
			}{},
			expectedErr: ErrInvalidFieldType,
		},
		{
			in: struct {
				Field string `validate:"regexp:(-"`
			}{},
			expectedErr: ErrInvalidRegexp,
		},
		{
			in: struct {
				Field []byte `validate:"max:11"`
			}{},
			expectedErr: ErrInvalidFieldType,
		},
		{
			in: App{Version: "1"},
			expectedErr: ValidationErrors{{
				Field: "Version",
				Value: "1",
				Err:   ErrValidationLen,
			}},
		},
		{
			in: Response{Code: 100, Body: "tt"},
			expectedErr: ValidationErrors{{
				Field: "Code",
				Value: 100,
				Err:   ErrValidationIn,
			}},
		},
		{
			in: User{
				ID:     "123",
				Age:    11,
				Email:  "test/mail.com",
				Role:   "test",
				Phones: []string{"123", "1234"},
			},
			expectedErr: ValidationErrors{
				{
					Field: "ID",
					Value: "123",
					Err:   ErrValidationLen,
				},
				{
					Field: "Age",
					Value: 11,
					Err:   ErrValidationMin,
				},
				{
					Field: "Email",
					Value: "test@mailcom",
					Err:   ErrValidationRegexp,
				},
				{
					Field: "Role",
					Value: "test",
					Err:   ErrValidationIn,
				},
				{
					Field: "Phones",
					Value: "123",
					Err:   ErrValidationLen,
				},
				{
					Field: "Phones",
					Value: "1234",
					Err:   ErrValidationLen,
				},
			},
		},
		{
			in: User{
				ID:     "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
				Age:    30,
				Email:  "test@mail.com",
				Role:   "admin",
				Phones: []string{"11111111111"},
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			var vErr ValidationErrors
			var expErr ValidationErrors

			if errors.As(err, &vErr) {
				if errors.As(tt.expectedErr, &expErr) {
					for i, verr := range vErr {
						assert.ErrorIs(t, verr.Err, expErr[i].Err)
					}
				}
			} else {
				assert.ErrorIs(t, err, tt.expectedErr)
			}

			if tt.expectedErr == nil {
				assert.NoError(t, err)
			}

			_ = tt
		})
	}
}
