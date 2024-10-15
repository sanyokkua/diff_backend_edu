package services

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const testJwtKey = "testSecretKey" // Example secret key for testing

func TestJwtService(t *testing.T) {
	service := NewJwtService([]byte(testJwtKey))

	t.Run("GenerateJwtToken_Success", func(t *testing.T) {
		username := "testuser"
		token, err := service.GenerateJwtToken(username)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Extract claims to check if the token contains the expected username
		claims, err := service.ExtractClaims(token)
		assert.NoError(t, err)

		claimUsername, _ := claims.GetSubject()
		assert.Equal(t, username, claimUsername)
	})

	t.Run("ExtractClaims_Success", func(t *testing.T) {
		username := "testuser"
		token, _ := service.GenerateJwtToken(username)

		claims, err := service.ExtractClaims(token)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
	})

	t.Run("ExtractClaims_InvalidToken", func(t *testing.T) {
		_, err := service.ExtractClaims("invalid.token.string")
		assert.Error(t, err)
	})

	t.Run("IsTokenExpired_Success", func(t *testing.T) {
		username := "testuser"
		token, _ := service.GenerateJwtToken(username)
		claims, _ := service.ExtractClaims(token)

		expired := service.IsTokenExpired(claims)
		assert.False(t, expired)

		// Simulate expiration by manually setting the expiration claim to the past
		claims.(jwt.MapClaims)["exp"] = time.Now().Add(-time.Minute).Unix()
		expired = service.IsTokenExpired(claims)
		assert.True(t, expired)
	})

	t.Run("ValidateToken_Success", func(t *testing.T) {
		username := "testuser"
		token, _ := service.GenerateJwtToken(username)

		valid := service.ValidateToken(token, username)
		assert.True(t, valid)
	})

	t.Run("ValidateToken_InvalidUsername", func(t *testing.T) {
		username := "testuser"
		token, _ := service.GenerateJwtToken(username)

		valid := service.ValidateToken(token, "wronguser")
		assert.False(t, valid)
	})

	t.Run("ValidateToken_InvalidToken", func(t *testing.T) {
		valid := service.ValidateToken("invalid.token.string", "testuser")
		assert.False(t, valid)
	})

	t.Run("ValidateToken_EmptyTokenOrUsername", func(t *testing.T) {
		valid := service.ValidateToken("", "testuser")
		assert.False(t, valid)

		valid = service.ValidateToken("some.token", "")
		assert.False(t, valid)
	})
}
