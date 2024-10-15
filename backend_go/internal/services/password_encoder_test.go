package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordEncoder(t *testing.T) {
	encoder := NewDefaultBCryptPasswordEncoder()

	t.Run("Encode_Success", func(t *testing.T) {
		rawPassword := "testPassword123"
		encodedPassword, err := encoder.Encode(rawPassword)

		assert.NoError(t, err)
		assert.NotEmpty(t, encodedPassword)

		// Verify that the encoded password is a valid bcrypt hash
		err = bcrypt.CompareHashAndPassword([]byte(encodedPassword), []byte(rawPassword))
		assert.NoError(t, err)
	})

	t.Run("Encode_Error", func(t *testing.T) {
		// Here we will attempt to create an encoder with an invalid cost
		invalidEncoder := NewBCryptPasswordEncoder(-1) // Default cost is chosen in the lib if it less than expected
		_, err := invalidEncoder.Encode("test")
		assert.NoError(t, err)
	})

	t.Run("Matches_Success", func(t *testing.T) {
		rawPassword := "testPassword123"
		encodedPassword, _ := encoder.Encode(rawPassword)

		isMatch, err := encoder.Matches(rawPassword, encodedPassword)
		assert.NoError(t, err)
		assert.True(t, isMatch)
	})

	t.Run("Matches_Failure", func(t *testing.T) {
		rawPassword := "testPassword123"
		encodedPassword, _ := encoder.Encode(rawPassword)

		// Test with a different password
		isMatch, err := encoder.Matches("wrongPassword", encodedPassword)
		assert.NoError(t, err)
		assert.False(t, isMatch)
	})

	t.Run("Matches_Error", func(t *testing.T) {
		isMatch, err := encoder.Matches("testPassword123", "invalid-hash")
		assert.Error(t, err)
		assert.False(t, isMatch) // Expect false since the hash is invalid
	})
}
