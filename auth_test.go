package twitter

import (
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

func TestLoginInput_Sanitize(t *testing.T) {
	input := LoginInput{
		Email:    " USER@gmail.com ",
		Password: "password",
	}

	input.Sanitize()
	require.Equal(t, input.Email, "user@gmail.com")
}

func TestLoginInput_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		input LoginInput
		err   error
	}{
		{
			name: "valid input",
			input: LoginInput{
				Email:    "userTest@gmail.com",
				Password: "password",
			},
			err: nil,
		}, {
			name: "invalid email",
			input: LoginInput{
				Email:    "user",
				Password: "password",
			},
			err: ErrValidation,
		}, {
			name: "empty password",
			input: LoginInput{
				Email:    "userTest@gmail.com",
				Password: "",
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
