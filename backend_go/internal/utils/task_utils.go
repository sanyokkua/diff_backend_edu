package utils

import (
	"errors"
	"github.com/rs/zerolog/log"
	"go_backend/internal/api"
	"go_backend/internal/apperrors"
	"go_backend/internal/dto"
	"go_backend/internal/models"
	"gorm.io/gorm"
)

// CheckTaskExistsForUser checks if a task with the given name exists for the user.
//
// Parameters:
//   - taskRepo: The repository to access task data.
//   - user: The user to check for the task's existence.
//   - taskName: The name of the task to check for.
//
// Returns:
//   - error: An error if the task already exists for the user, or if an error occurs during the check.
//     Returns nil if the task does not exist.
func CheckTaskExistsForUser(taskRepo api.TaskRepository, user *models.User, taskName string) error {
	log.Debug().Str("taskName", taskName).Msg("Checking if task exists for user")

	_, err := taskRepo.FindByUserAndName(user, taskName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Debug().Str("taskName", taskName).Msg("Task not found for user")
			return nil
		}
		log.Error().Err(err).Str("taskName", taskName).Msg("Error finding task for user")
		return err
	}

	log.Warn().Str("taskName", taskName).Msg("Task already exists for user")
	return apperrors.NewTaskAlreadyExistsError("Task with the name '" + taskName + "' already exists for the user")
}

// ValidateTaskCreation validates the task creation DTO.
//
// Parameters:
//   - taskCreationDTO: The task creation data transfer object to validate.
//
// Returns:
//   - error: An error if the validation fails. Returns nil if the validation passes.
func ValidateTaskCreation(taskCreationDTO *dto.TaskCreationDTO) error {
	log.Debug().Msg("Validating TaskCreationDTO")

	if taskCreationDTO == nil {
		log.Error().Msg("TaskCreationDTO is nil")
		return apperrors.NewIllegalArgumentError("TaskCreationDTO is nil")
	}

	if taskCreationDTO.Name == "" {
		log.Error().Msg("Task name cannot be empty")
		return apperrors.NewIllegalArgumentError("Task name cannot be empty")
	}

	if taskCreationDTO.Description == "" {
		log.Error().Msg("Task description cannot be empty")
		return apperrors.NewIllegalArgumentError("Task description cannot be empty")
	}

	log.Debug().Msg("TaskCreationDTO is valid")
	return nil
}

// ValidateTaskUpdateDTO validates the task update DTO.
//
// Parameters:
//   - taskUpdateDTO: The task update data transfer object to validate.
//
// Returns:
//   - error: An error if the validation fails. Returns nil if the validation passes.
func ValidateTaskUpdateDTO(taskUpdateDTO *dto.TaskUpdateDTO) error {
	log.Debug().Msg("Validating TaskUpdateDTO")

	if taskUpdateDTO == nil {
		log.Error().Msg("TaskUpdateDTO is nil")
		return apperrors.NewIllegalArgumentError("TaskUpdateDTO is nil")
	}

	if taskUpdateDTO.Name == "" {
		log.Error().Msg("Task name cannot be empty")
		return apperrors.NewIllegalArgumentError("Task name cannot be empty")
	}

	if taskUpdateDTO.Description == "" {
		log.Error().Msg("Task description cannot be empty")
		return apperrors.NewIllegalArgumentError("Task description cannot be empty")
	}

	log.Debug().Msg("TaskUpdateDTO is valid")
	return nil
}
