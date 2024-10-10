package repositories

import (
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"go_backend/models"
	"gorm.io/gorm"
	"testing"
)

func TestCreateTask(t *testing.T) {
	// Setup test container
	postgresContainer, db, err := GetDBForTests()
	if err != nil {
		log.Err(err).Msg("Failed to set up test DB")
		t.FailNow()
	}
	defer cleanUpDB(postgresContainer, db) // Clean up after test run

	userRepo := NewUserRepository(db)
	repository := NewTaskRepository(db)

	testUser, err := userRepo.CreateUser(&models.User{Email: "user@email.com", PasswordHash: "SecurePassword"})
	if err != nil {
		t.Fail()
	}

	t.Run("CreateTask_Success", func(t *testing.T) {
		user := testUser
		task := &models.Task{
			Name:        "Test Task",
			Description: "This is a test task",
			UserID:      user.UserID,
		}

		createdTask, err := repository.CreateTask(task)

		assert.NoError(t, err)
		assert.NotNil(t, createdTask)
		assert.Equal(t, task.Name, createdTask.Name)
	})

	t.Run("CreateTask_NilTask", func(t *testing.T) {
		_, err := repository.CreateTask(nil)

		assert.Error(t, err)
		assert.Equal(t, "Task model is nil", err.Error())
	})
}

func TestFindByTaskID(t *testing.T) {
	postgresContainer, db, err := GetDBForTests()
	if err != nil {
		log.Err(err).Msg("Failed to set up test DB")
		t.FailNow()
	}
	defer cleanUpDB(postgresContainer, db) // Clean up after test run

	userRepo := NewUserRepository(db)
	repository := NewTaskRepository(db)

	testUser, err := userRepo.CreateUser(&models.User{Email: "user@email.com", PasswordHash: "SecurePassword"})
	if err != nil {
		t.Fail()
	}

	t.Run("FindByTaskID_Success", func(t *testing.T) {
		user := testUser
		task := &models.Task{
			Name:        "Test Task",
			Description: "This is a test task",
			UserID:      user.UserID,
		}
		createdTask, _ := repository.CreateTask(task)

		retrievedTask, err := repository.FindByTaskID(createdTask.TaskID)

		assert.NoError(t, err)
		assert.NotNil(t, retrievedTask)
		assert.Equal(t, createdTask.TaskID, retrievedTask.TaskID)
	})

	t.Run("FindByTaskID_InvalidID", func(t *testing.T) {
		_, err := repository.FindByTaskID(-1)

		assert.Error(t, err)
		assert.Equal(t, "ID is not valid", err.Error())
	})

	t.Run("FindByTaskID_NotFound", func(t *testing.T) {
		_, err := repository.FindByTaskID(9999)

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestFindByUserAndName(t *testing.T) {
	postgresContainer, db, err := GetDBForTests()
	if err != nil {
		log.Err(err).Msg("Failed to set up test DB")
		t.FailNow()
	}
	defer cleanUpDB(postgresContainer, db) // Clean up after test run

	userRepo := NewUserRepository(db)
	repository := NewTaskRepository(db)

	testUser, err := userRepo.CreateUser(&models.User{Email: "user@email.com", PasswordHash: "SecurePassword"})
	if err != nil {
		t.Fail()
	}

	t.Run("FindByUserAndName_Success", func(t *testing.T) {
		user := testUser
		task := &models.Task{
			Name:        "Unique Task Name",
			Description: "This is a unique task",
			UserID:      user.UserID,
		}
		_, _ = repository.CreateTask(task)

		retrievedTask, err := repository.FindByUserAndName(user, task.Name)

		assert.NoError(t, err)
		assert.NotNil(t, retrievedTask)
		assert.Equal(t, task.Name, retrievedTask.Name)
	})

	t.Run("FindByUserAndName_InvalidParams", func(t *testing.T) {
		_, err := repository.FindByUserAndName(nil, "Some Task Name")

		assert.Error(t, err)
		assert.Equal(t, "Invalid params", err.Error())

		_, err = repository.FindByUserAndName(&models.User{UserID: 1}, "")

		assert.Error(t, err)
		assert.Equal(t, "Invalid params", err.Error())
	})

	t.Run("FindByUserAndName_NotFound", func(t *testing.T) {
		user := &models.User{UserID: 1} // Assuming a user with ID 1 exists
		_, err := repository.FindByUserAndName(user, "Nonexistent Task")

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestFindAllByUser(t *testing.T) {
	postgresContainer, db, err := GetDBForTests()
	if err != nil {
		log.Err(err).Msg("Failed to set up test DB")
		t.FailNow()
	}
	defer cleanUpDB(postgresContainer, db) // Clean up after test run

	userRepo := NewUserRepository(db)
	repository := NewTaskRepository(db)

	testUser, err := userRepo.CreateUser(&models.User{Email: "user@email.com", PasswordHash: "SecurePassword"})
	if err != nil {
		t.Fail()
	}

	t.Run("FindAllByUser_Success", func(t *testing.T) {
		user := testUser
		task1 := &models.Task{Name: "Task 1", Description: "First task", UserID: user.UserID}
		task2 := &models.Task{Name: "Task 2", Description: "Second task", UserID: user.UserID}
		_, _ = repository.CreateTask(task1)
		_, _ = repository.CreateTask(task2)

		tasks, err := repository.FindAllByUser(user)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(tasks))
	})

	t.Run("FindAllByUser_NilUser", func(t *testing.T) {
		_, err := repository.FindAllByUser(nil)

		assert.Error(t, err)
		assert.Equal(t, "User is nil", err.Error())
	})
}

func TestUpdateTask(t *testing.T) {
	postgresContainer, db, err := GetDBForTests()
	if err != nil {
		log.Err(err).Msg("Failed to set up test DB")
		t.FailNow()
	}
	defer cleanUpDB(postgresContainer, db) // Clean up after test run

	userRepo := NewUserRepository(db)
	repository := NewTaskRepository(db)

	testUser, err := userRepo.CreateUser(&models.User{Email: "user@email.com", PasswordHash: "SecurePassword"})
	if err != nil {
		t.Fail()
	}

	t.Run("UpdateTask_Success", func(t *testing.T) {
		user := testUser
		task := &models.Task{Name: "Update Task", Description: "Initial description", UserID: user.UserID}
		createdTask, _ := repository.CreateTask(task)

		createdTask.Description = "Updated description"
		updatedTask, err := repository.UpdateTask(createdTask)

		assert.NoError(t, err)
		assert.Equal(t, "Updated description", updatedTask.Description)
	})

	t.Run("UpdateTask_NilTask", func(t *testing.T) {
		_, err := repository.UpdateTask(nil)

		assert.Error(t, err)
		assert.Equal(t, "Task model is nil", err.Error())
	})
}

func TestDeleteTask(t *testing.T) {
	postgresContainer, db, err := GetDBForTests()
	if err != nil {
		log.Err(err).Msg("Failed to set up test DB")
		t.FailNow()
	}
	defer cleanUpDB(postgresContainer, db) // Clean up after test run

	userRepo := NewUserRepository(db)
	repository := NewTaskRepository(db)

	testUser, err := userRepo.CreateUser(&models.User{Email: "user@email.com", PasswordHash: "SecurePassword"})
	if err != nil {
		t.Fail()
	}

	t.Run("DeleteTask_Success", func(t *testing.T) {
		user := testUser
		task := &models.Task{Name: "Delete Task", Description: "This task will be deleted", UserID: user.UserID}
		createdTask, _ := repository.CreateTask(task)

		err := repository.DeleteTask(createdTask)
		assert.NoError(t, err)

		// Verify the task is deleted
		_, err = repository.FindByTaskID(createdTask.TaskID)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteTask_NilTask", func(t *testing.T) {
		err := repository.DeleteTask(nil)

		assert.Error(t, err)
		assert.Equal(t, "Task model is nil", err.Error())
	})
}

func TestFindByUserAndTaskID(t *testing.T) {
	postgresContainer, db, err := GetDBForTests()
	if err != nil {
		log.Err(err).Msg("Failed to set up test DB")
		t.FailNow()
	}
	defer cleanUpDB(postgresContainer, db) // Clean up after test run

	userRepo := NewUserRepository(db)
	repository := NewTaskRepository(db)

	testUser, err := userRepo.CreateUser(&models.User{Email: "user@email.com", PasswordHash: "SecurePassword"})
	if err != nil {
		t.Fail()
	}

	t.Run("FindByUserAndTaskID_Success", func(t *testing.T) {
		user := testUser
		task := &models.Task{
			Name:        "Task for User",
			Description: "This task belongs to a user",
			UserID:      user.UserID,
		}
		createdTask, _ := repository.CreateTask(task)

		retrievedTask, err := repository.FindByUserAndTaskID(user, createdTask.TaskID)

		assert.NoError(t, err)
		assert.NotNil(t, retrievedTask)
		assert.Equal(t, createdTask.TaskID, retrievedTask.TaskID)
		assert.Equal(t, user.UserID, retrievedTask.UserID)
	})

	t.Run("FindByUserAndTaskID_InvalidParams", func(t *testing.T) {
		_, err := repository.FindByUserAndTaskID(nil, 1)

		assert.Error(t, err)
		assert.Equal(t, "Invalid params", err.Error())

		user := &models.User{UserID: 1} // Assuming a user with ID 1 exists
		_, err = repository.FindByUserAndTaskID(user, -1)

		assert.Error(t, err)
		assert.Equal(t, "Invalid params", err.Error())
	})

	t.Run("FindByUserAndTaskID_NotFound", func(t *testing.T) {
		user := &models.User{UserID: 1}                      // Assuming a user with ID 1 exists
		_, err := repository.FindByUserAndTaskID(user, 9999) // Non-existent task ID

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}
