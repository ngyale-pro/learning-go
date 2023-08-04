package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	testCases := []struct {
		name          string
		checkResponse func(t *testing.T, password string)
	}{
		{
			name: "OK",
			checkResponse: func(t *testing.T, password string) {
				hashedPassword, err := HashPassword(password)
				require.NoError(t, err)
				require.NotEmpty(t, hashedPassword)
				err = CheckPassword(hashedPassword, password)
				require.NoError(t, err)
			},
		},
		{
			name: "Wrong Password",
			checkResponse: func(t *testing.T, password string) {
				hashedPassword, err := HashPassword(password)
				require.NoError(t, err)
				require.NotEmpty(t, hashedPassword)
				wrongPassword := RandomString(10)
				err = CheckPassword(hashedPassword, wrongPassword)
				require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
			},
		},
		{
			name: "Same Password Gives Different Hash(salt)",
			checkResponse: func(t *testing.T, password string) {
				hashedPassword1, err := HashPassword(password)
				require.NoError(t, err)
				require.NotEmpty(t, hashedPassword1)

				hashedPassword2, err := HashPassword(password)
				require.NoError(t, err)
				require.NotEmpty(t, hashedPassword2)
				require.NotEqual(t, hashedPassword1, hashedPassword2)
			},
		},
	}

	for _, testCase := range testCases {
		password := RandomString(10)
		t.Run(testCase.name, func(t *testing.T) {
			testCase.checkResponse(t, password)
		})
	}

}
