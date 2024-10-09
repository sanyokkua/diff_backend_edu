package repositories

import (
	"github.com/rs/zerolog/log"
	"go_backend/api"
	"go_backend/models"
	"gorm.io/gorm"
)

// userRepository is a private struct that implements the UserRepository interface.
// It provides methods for interacting with users in the database.
type userRepository struct {
	db *gorm.DB // Database connection used to interact with the users table.
}

// NewUserRepository creates a new instance of UserRepository.
// It initializes a userRepository with the given database connection.
//
// Parameters:
//   - db: A pointer to the gorm.DB instance used for database operations.
//
// Returns:
//   - api.UserRepository: A new instance of UserRepository.
func NewUserRepository(db *gorm.DB) api.UserRepository {
	log.Debug().Msg("Creating new UserRepository")
	return &userRepository{db: db}
}

// CreateUser inserts a new user into the database.
//
// Parameters:
//   - user: A pointer to the User model to be created.
//
// Returns:
//   - *models.User: A pointer to the created user.
//   - error: An error if the creation fails.
func (r *userRepository) CreateUser(user *models.User) (*models.User, error) {
	log.Debug().Str("email", user.Email).Msg("Creating user")
	if err := r.db.Create(user).Error; err != nil {
		log.Error().Err(err).Str("email", user.Email).Msg("Error creating user")
		return nil, err
	}
	log.Debug().Str("email", user.Email).Msg("User created successfully")
	return user, nil
}

// GetUserByID retrieves a user by their ID.
//
// Parameters:
//   - id: The unique identifier of the user to be retrieved.
//
// Returns:
//   - *models.User: A pointer to the user with the specified ID.
//   - error: An error if the retrieval fails.
func (r *userRepository) GetUserByID(id int64) (*models.User, error) {
	log.Debug().Int64("id", id).Msg("Retrieving user by ID")
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		log.Error().Err(err).Int64("id", id).Msg("Error retrieving user by ID")
		return nil, err
	}
	log.Debug().Int64("id", id).Msg("User retrieved successfully")
	return &user, nil
}

// UpdateUser updates an existing user in the database.
//
// Parameters:
//   - user: A pointer to the User model to be updated.
//
// Returns:
//   - *models.User: A pointer to the updated user.
//   - error: An error if the update fails.
func (r *userRepository) UpdateUser(user *models.User) (*models.User, error) {
	log.Debug().Str("email", user.Email).Msg("Updating user")
	if err := r.db.Save(user).Error; err != nil {
		log.Error().Err(err).Str("email", user.Email).Msg("Error updating user")
		return nil, err
	}
	log.Debug().Str("email", user.Email).Msg("User updated successfully")
	return user, nil
}

// DeleteUser removes a user from the database by their ID.
//
// Parameters:
//   - id: The unique identifier of the user to be deleted.
//
// Returns:
//   - error: An error if the deletion fails.
func (r *userRepository) DeleteUser(id int64) error {
	log.Debug().Int64("id", id).Msg("Deleting user")
	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		log.Error().Err(err).Int64("id", id).Msg("Error deleting user")
		return err
	}
	log.Debug().Int64("id", id).Msg("User deleted successfully")
	return nil
}

// GetUserByEmail retrieves a user by their email address.
//
// Parameters:
//   - email: The email address of the user to be retrieved.
//
// Returns:
//   - *models.User: A pointer to the user with the specified email address.
//   - error: An error if the retrieval fails.
func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	log.Debug().Str("email", email).Msg("Retrieving user by email")
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		log.Error().Err(err).Str("email", email).Msg("Error retrieving user by email")
		return nil, err
	}
	log.Debug().Str("email", email).Msg("User retrieved successfully")
	return &user, nil
}
