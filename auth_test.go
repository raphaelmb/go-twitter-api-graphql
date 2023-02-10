package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterInput_Sanitize(t *testing.T) {
	input := RegisterInput{
		Username:        " username ",
		Email:           " USER@gmail.com ",
		Password:        "password",
		ConfirmPassword: "password",
	}

	input.Sanitize()

	fmt.Println(input)

	require.Equal(t, input.Username, "username")
	require.Equal(t, input.Email, "user@gmail.com")
}

func TestRegisterInput_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		input RegisterInput
		err   error
	}{
		{
			name: "valid input",
			input: RegisterInput{
				Username:        "userTest",
				Email:           "userTest@gmail.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: nil,
		}, {
			name: "invalid email",
			input: RegisterInput{
				Username:        "user",
				Email:           "user",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: ErrValidation,
		}, {
			name: "username too short",
			input: RegisterInput{
				Username:        "u",
				Email:           "userTest@gmail.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: ErrValidation,
		}, {
			name: "password too short",
			input: RegisterInput{
				Username:        "userTest",
				Email:           "userTest@gmail.com",
				Password:        "pass",
				ConfirmPassword: "pass",
			},
			err: ErrValidation,
		}, {
			name: "password does not match confirm password",
			input: RegisterInput{
				Username:        "userTest",
				Email:           "userTest@gmail.com",
				Password:        "password1",
				ConfirmPassword: "password2",
			},
			err: ErrValidation,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
