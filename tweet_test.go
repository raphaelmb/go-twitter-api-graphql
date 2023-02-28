package twitter

import (
	"testing"

	"github.com/raphaelmb/go-twitter-api-graphql/faker"
	"github.com/stretchr/testify/require"
)

func TestCreateTweetInput_Sanitize(t *testing.T) {
	input := CreateTweetInput{
		Body: " hello   ",
	}

	expected := CreateTweetInput{
		Body: "hello",
	}

	input.Sanitize()

	require.Equal(t, input, expected)
}

func TestCreateTweetInput_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		input CreateTweetInput
		err   error
	}{
		{
			name: "valid",
			input: CreateTweetInput{
				Body: "hello",
			},
			err: nil,
		}, {
			name: "tweet not long enough",
			input: CreateTweetInput{
				Body: "h",
			},
			err: ErrValidation,
		}, {
			name: "tweet too long",
			input: CreateTweetInput{
				Body: faker.RandStr(300),
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
