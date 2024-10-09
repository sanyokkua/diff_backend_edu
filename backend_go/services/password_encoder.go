package services

import (
	"errors"
	"github.com/rs/zerolog/log"
	"go_backend/api"
	"golang.org/x/crypto/bcrypt"
)

// passwordEncoder implements the PasswordEncoder interface and provides methods for encoding and matching passwords.
type passwordEncoder struct {
	cost int // The cost factor for bcrypt hashing
}

// NewDefaultBCryptPasswordEncoder creates a new password encoder with the default bcrypt cost.
//
// Returns:
//   - api.PasswordEncoder: A new instance of PasswordEncoder with default bcrypt cost.
func NewDefaultBCryptPasswordEncoder() api.PasswordEncoder {
	log.Debug().Msg("Creating default BCryptPasswordEncoder")
	return NewBCryptPasswordEncoder(bcrypt.DefaultCost)
}

// NewBCryptPasswordEncoder creates a new password encoder with the specified bcrypt cost.
//
// Parameters:
//   - cost: An integer representing the cost factor for bcrypt hashing. Higher costs increase the time required to hash the password.
//
// Returns:
//   - api.PasswordEncoder: A new instance of PasswordEncoder with the specified bcrypt cost.
func NewBCryptPasswordEncoder(cost int) api.PasswordEncoder {
	log.Debug().Int("cost", cost).Msg("Creating BCryptPasswordEncoder with specified cost")
	return &passwordEncoder{cost: cost}
}

// Matches checks if the raw password matches the encoded password.
//
// Parameters:
//   - rawPassword: A string representing the plaintext password to verify.
//   - encodedPassword: A string representing the previously encoded password (hash).
//
// Returns:
//   - bool: True if the passwords match, otherwise false.
//   - error: An error if the comparison fails due to an unexpected error.
func (p *passwordEncoder) Matches(rawPassword, encodedPassword string) (bool, error) {
	log.Debug().Msg("Checking if passwords match")

	err := bcrypt.CompareHashAndPassword([]byte(encodedPassword), []byte(rawPassword))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		log.Debug().Msg("Passwords do not match")
		return false, nil // Return false if passwords do not match
	}

	if err != nil {
		log.Error().Err(err).Msg("Error comparing hash and password")
		return false, err // Return the error if comparison fails
	}

	log.Debug().Msg("Passwords match")
	return true, nil // Return true if passwords match
}

// Encode hashes the password using bcrypt.
//
// Parameters:
//   - password: A string representing the plaintext password to be hashed.
//
// Returns:
//   - string: The hashed password as a string.
//   - error: An error if the hashing process fails.
func (p *passwordEncoder) Encode(password string) (string, error) {
	log.Debug().Msg("Encoding password")

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	if err != nil {
		log.Error().Err(err).Msg("Error generating password hash")
		return "", err // Return the error if hashing fails
	}

	hashedPassword := string(hashedBytes)
	log.Debug().Msg("Password encoded successfully")
	return hashedPassword, nil // Return the hashed password
}
