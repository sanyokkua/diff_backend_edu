package services

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"go_backend/internal/api"
	"time"
)

const JwtExpMinutes = 15 // JWT expiration time in minutes

// jwtService implements the JwtService interface and provides methods for managing JWT tokens.
type jwtService struct {
	jwtMainKey []byte // The main key used for signing JWT tokens.
}

// NewJwtService creates a new JwtService with the given main key.
//
// Parameters:
//   - jwtMainKey: A byte slice representing the key used for signing tokens.
//
// Returns:
//   - api.JwtService: A new instance of JwtService.
func NewJwtService(jwtMainKey []byte) api.JwtService {
	log.Debug().Msg("Creating new JwtService")
	return &jwtService{jwtMainKey: jwtMainKey}
}

// ExtractClaims extracts claims from a JWT token after verifying its signature.
//
// Parameters:
//   - token: A string representing the JWT token from which to extract claims.
//
// Returns:
//   - jwt.Claims: The extracted claims from the token.
//   - error: An error if the extraction fails.
func (s *jwtService) ExtractClaims(token string) (jwt.Claims, error) {
	log.Debug().Str("token", token).Msg("Extracting claims from token")

	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtMainKey, nil // Provide the signing key for verification
	}, jwt.WithoutClaimsValidation()) // Allow claims extraction without validating expiration

	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		log.Error().Err(err).Msg("Error parsing token")
		return nil, err
	}

	log.Debug().Msg("Claims extracted successfully")
	return parsedToken.Claims, nil
}

// IsTokenExpired checks if the JWT token has expired based on the claims.
//
// Parameters:
//   - claims: The claims extracted from the JWT token.
//
// Returns:
//   - bool: True if the token is expired, otherwise false.
func (s *jwtService) IsTokenExpired(claims jwt.Claims) bool {
	log.Debug().Msg("Checking if token is expired")

	expTime, err := claims.GetExpirationTime()
	if err != nil {
		log.Warn().Msg("Token expiration date is missing.")
		return true // Consider token expired if expiration time is missing
	}

	expirationTime := time.Unix(expTime.UnixNano(), 0)
	isExpired := expirationTime.Before(time.Now())
	log.Debug().Bool("isExpired", isExpired).Msg("Token expiration check result")
	return isExpired
}

// ValidateToken checks if the provided JWT token is valid and if the username matches the claims.
//
// Parameters:
//   - token: A string representing the JWT token to validate.
//   - username: A string representing the username to compare against the token claims.
//
// Returns:
//   - bool: True if the token is valid and the username matches, otherwise false.
func (s *jwtService) ValidateToken(token, username string) bool {
	log.Debug().Str("token", token).Str("username", username).Msg("Validating token")

	if token == "" || username == "" {
		log.Warn().Msg("Token or username is missing during validation.")
		return false
	}

	claims, err := s.ExtractClaims(token)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to extract claims from token")
		return false
	}

	claimUsername, subjErr := claims.GetSubject()
	if subjErr != nil {
		log.Warn().Msg("Username claim is missing in the token.")
		return false
	}

	isUsernameValid := username == claimUsername
	isTokenValid := !s.IsTokenExpired(claims)

	log.Debug().Bool("isUsernameValid", isUsernameValid).Bool("isTokenValid", isTokenValid).Msg("Token validation result")
	return isUsernameValid && isTokenValid
}

// GenerateJwtToken generates a JWT token for the specified username.
//
// Parameters:
//   - username: A string representing the username for which to generate the token.
//
// Returns:
//   - string: The generated JWT token as a string.
//   - error: An error if token generation fails.
func (s *jwtService) GenerateJwtToken(username string) (string, error) {
	log.Debug().Str("username", username).Msg("Generating JWT token")

	now := time.Now()
	exp := now.Add(time.Minute * JwtExpMinutes) // Set token expiration time

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,   // Subject of the token (the username)
		"iat": now.Unix(), // Issued at time
		"exp": exp.Unix(), // Expiration time
	})

	tokenString, err := token.SignedString(s.jwtMainKey)
	if err != nil {
		log.Error().Err(err).Msg("Error generating JWT token")
		return "", err
	}

	log.Debug().Str("username", username).Time("expiration", exp).Msg("Generated JWT token successfully")
	return tokenString, nil
}
