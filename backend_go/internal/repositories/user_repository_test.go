package repositories

import (
	"go_backend/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	_, db, cleanup := setupTest(t)
	defer cleanup()

	repository := NewUserRepository(db)

	t.Run("CreateUser_Success", func(t *testing.T) {
		user := &models.User{
			Email:        "test@example.com",
			PasswordHash: "password_hash",
		}

		createdUser, err := repository.CreateUser(user)

		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, user.Email, createdUser.Email)
	})

	t.Run("CreateUser_NilUser", func(t *testing.T) {
		_, err := repository.CreateUser(nil)

		assert.Error(t, err)
		assert.Equal(t, "Passed User model is nil", err.Error())
	})
}

func TestGetUserByID(t *testing.T) {
	_, db, cleanup := setupTest(t)
	defer cleanup()

	repository := NewUserRepository(db)

	t.Run("GetUserByID_Success", func(t *testing.T) {
		user := &models.User{
			Email:        "test2@example.com",
			PasswordHash: "password_hash2",
		}
		createdUser, _ := repository.CreateUser(user)

		retrievedUser, err := repository.GetUserByID(createdUser.UserID)

		assert.NoError(t, err)
		assert.NotNil(t, retrievedUser)
		assert.Equal(t, createdUser.UserID, retrievedUser.UserID)
	})

	t.Run("GetUserByID_InvalidID", func(t *testing.T) {
		_, err := repository.GetUserByID(-1)

		assert.Error(t, err)
		assert.Equal(t, "Passed id is invalid", err.Error())
	})

	t.Run("GetUserByID_NotFound", func(t *testing.T) {
		_, err := repository.GetUserByID(9999)

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestUpdateUser(t *testing.T) {
	_, db, cleanup := setupTest(t)
	defer cleanup()

	repository := NewUserRepository(db)

	t.Run("UpdateUser_Success", func(t *testing.T) {
		user := &models.User{
			Email:        "update_test@example.com",
			PasswordHash: "initial_hash",
		}
		createdUser, _ := repository.CreateUser(user)

		// Update the password hash
		createdUser.PasswordHash = "updated_hash"
		updatedUser, err := repository.UpdateUser(createdUser)

		assert.NoError(t, err)
		assert.Equal(t, "updated_hash", updatedUser.PasswordHash)
	})

	t.Run("UpdateUser_NilUser", func(t *testing.T) {
		_, err := repository.UpdateUser(nil)

		assert.Error(t, err)
		assert.Equal(t, "Passed User model is nil", err.Error())
	})
}

func TestDeleteUser(t *testing.T) {
	_, db, cleanup := setupTest(t)
	defer cleanup()

	repository := NewUserRepository(db)

	t.Run("DeleteUser_Success", func(t *testing.T) {
		user := &models.User{
			Email:        "delete_test@example.com",
			PasswordHash: "password_hash",
		}
		createdUser, _ := repository.CreateUser(user)

		err := repository.DeleteUser(createdUser.UserID)
		assert.NoError(t, err)

		// Verify the user is deleted
		_, err = repository.GetUserByID(createdUser.UserID)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteUser_InvalidID", func(t *testing.T) {
		err := repository.DeleteUser(-1)

		assert.Error(t, err)
		assert.Equal(t, "Passed id is invalid", err.Error())
	})
}

func TestGetUserByEmail(t *testing.T) {
	_, db, cleanup := setupTest(t)
	defer cleanup()

	repository := NewUserRepository(db)

	t.Run("GetUserByEmail_Success", func(t *testing.T) {
		user := &models.User{
			Email:        "email_test@example.com",
			PasswordHash: "password_hash",
		}
		_, err := repository.CreateUser(user)
		if err != nil {
			t.Fail()
		}

		retrievedUser, err := repository.GetUserByEmail(user.Email)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedUser)
		assert.Equal(t, user.Email, retrievedUser.Email)
	})

	t.Run("GetUserByEmail_EmptyEmail", func(t *testing.T) {
		_, err := repository.GetUserByEmail("")

		assert.Error(t, err)
		assert.Equal(t, "Email can't be empty string", err.Error())
	})

	t.Run("GetUserByEmail_NotFound", func(t *testing.T) {
		_, err := repository.GetUserByEmail("non_existent@example.com")

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}
