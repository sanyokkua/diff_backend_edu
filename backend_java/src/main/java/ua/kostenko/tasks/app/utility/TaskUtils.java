package ua.kostenko.tasks.app.utility;

import lombok.AccessLevel;
import lombok.NoArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import ua.kostenko.tasks.app.dto.task.TaskCreationDTO;
import ua.kostenko.tasks.app.dto.task.TaskUpdateDTO;
import ua.kostenko.tasks.app.entity.Task;
import ua.kostenko.tasks.app.entity.User;
import ua.kostenko.tasks.app.exception.TaskAlreadyExistsException;
import ua.kostenko.tasks.app.repository.TaskRepository;

import java.util.Objects;
import java.util.Optional;

/**
 * Utility class for task-related operations such as checking task existence and validating task creation/update data.
 *
 * <p>This class is designed as a singleton and cannot be instantiated due to the {@link NoArgsConstructor} with
 * {@link AccessLevel#PRIVATE}. It contains static methods to assist in task management, particularly in validation
 * and duplication checks.</p>
 */
@Slf4j
@NoArgsConstructor(access = AccessLevel.PRIVATE)
public class TaskUtils {

    /**
     * Checks if a task with the given name already exists for the specified user.
     *
     * <p>This method checks the task repository to see if a task with the same name already exists for a given user.
     * If such a task exists, it throws a {@link TaskAlreadyExistsException}.</p>
     *
     * @param taskRepository the repository to query for task existence
     * @param user           the user for whom the task is being created
     * @param taskName       the name of the task to check for duplicates
     *
     * @throws TaskAlreadyExistsException if a task with the given name already exists for the user
     */
    public static void checkTaskExistsForUser(TaskRepository taskRepository, User user, String taskName) {
        log.debug("Checking if task '{}' exists for user with ID: {}", taskName, user.getUserId());
        Optional<Task> existingTask = taskRepository.findByUserAndName(user, taskName);

        if (existingTask.isPresent()) {
            log.warn("Task '{}' already exists for user with ID: {}", taskName, user.getUserId());
            throw new TaskAlreadyExistsException("Task with the name '" + taskName + "' already exists for the user");
        }
        log.info("Task '{}' is available for creation for user with ID: {}", taskName, user.getUserId());
    }

    /**
     * Validates the task creation DTO.
     *
     * <p>This method validates the task creation data, ensuring that both the task name and description are not
     * null or empty. If any of the fields are invalid, it throws an {@link IllegalArgumentException}.</p>
     *
     * @param taskCreationDTO the task creation DTO containing the task data
     *
     * @throws IllegalArgumentException if any required fields (name, description) are invalid
     */
    public static void validateTaskCreation(TaskCreationDTO taskCreationDTO) {
        if (Objects.isNull(taskCreationDTO)) {
            log.error("TaskCreationDTO is null.");
            throw new IllegalArgumentException("TaskCreationDTO is null");
        }

        log.debug("Validating task creation data for task '{}'", taskCreationDTO.getName());

        if (taskCreationDTO.getName() == null || taskCreationDTO.getName().isBlank()) {
            log.warn("Task name is null or empty for task '{}'", taskCreationDTO.getName());
            throw new IllegalArgumentException("Task name cannot be null or empty");
        }

        if (taskCreationDTO.getDescription() == null || taskCreationDTO.getDescription().isBlank()) {
            log.warn("Task description is null or empty for task '{}'", taskCreationDTO.getName());
            throw new IllegalArgumentException("Task description cannot be null or empty");
        }

        log.info("Task creation data is valid for task '{}'", taskCreationDTO.getName());
    }

    /**
     * Validates the task update DTO.
     *
     * <p>This method ensures that the task update data is valid by checking that the task name and description
     * are not null or empty. It throws an {@link IllegalArgumentException} if any of the fields are invalid.</p>
     *
     * @param taskUpdateDTO the task update DTO containing the task data
     *
     * @throws IllegalArgumentException if any required fields (name, description) are invalid
     */
    public static void validateTaskUpdateDTO(TaskUpdateDTO taskUpdateDTO) {
        if (Objects.isNull(taskUpdateDTO)) {
            log.error("TaskUpdateDTO is null.");
            throw new IllegalArgumentException("TaskUpdateDTO is null");
        }

        log.debug("Validating task update data for task '{}'", taskUpdateDTO.getName());

        if (taskUpdateDTO.getName() == null || taskUpdateDTO.getName().isBlank()) {
            log.warn("Task name is null or empty for task '{}'", taskUpdateDTO.getName());
            throw new IllegalArgumentException("Task name cannot be null or empty");
        }

        if (taskUpdateDTO.getDescription() == null) {
            log.warn("Task description is null for task '{}'", taskUpdateDTO.getName());
            throw new IllegalArgumentException("Task description cannot be null");
        }

        log.info("Task update data is valid for task '{}'", taskUpdateDTO.getName());
    }
}

