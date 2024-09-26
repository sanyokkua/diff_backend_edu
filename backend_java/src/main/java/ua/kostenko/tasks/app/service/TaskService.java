package ua.kostenko.tasks.app.service;

import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.access.AccessDeniedException;
import org.springframework.stereotype.Service;
import ua.kostenko.tasks.app.dto.task.TaskCreationDTO;
import ua.kostenko.tasks.app.dto.task.TaskDto;
import ua.kostenko.tasks.app.dto.task.TaskUpdateDTO;
import ua.kostenko.tasks.app.entity.Task;
import ua.kostenko.tasks.app.entity.User;
import ua.kostenko.tasks.app.exception.TaskNotFoundException;
import ua.kostenko.tasks.app.repository.TaskRepository;
import ua.kostenko.tasks.app.repository.UserRepository;
import ua.kostenko.tasks.app.utility.TaskUtils;

import java.util.List;
import java.util.Optional;

/**
 * Service class for managing task-related operations such as creation, update, deletion, and retrieval of tasks.
 * <p>
 * This service ensures that only users who own the task can perform actions on it. It includes validation for
 * task creation and ensures tasks belong to the user performing the operation. All operations are logged for
 * traceability and debugging purposes.
 * </p>
 */
@Slf4j
@Service
@RequiredArgsConstructor
@Transactional
public class TaskService {

    private static final String USER_NOT_FOUND_EX = "User not found";
    private static final String TASK_NOT_FOUND_EX = "Task not found";
    private final TaskRepository taskRepository;
    private final UserRepository userRepository;

    /**
     * Creates a new task for the specified user.
     *
     * @param userId          the ID of the user creating the task
     * @param taskCreationDTO the DTO containing task creation details
     *
     * @return a {@link TaskDto} representing the newly created task
     *
     * @throws IllegalArgumentException if the user does not exist or if a task with the same name already exists
     */
    public TaskDto createTask(Long userId, TaskCreationDTO taskCreationDTO) {
        // Validate task data
        TaskUtils.validateTaskCreation(taskCreationDTO);
        log.info("Creating task '{}' for user with ID: {}", taskCreationDTO.getName(), userId);

        // Find the user
        User user = userRepository.findById(userId).orElseThrow(() -> new IllegalArgumentException(USER_NOT_FOUND_EX));

        // Check if a task with the same name exists for this user
        TaskUtils.checkTaskExistsForUser(taskRepository, user, taskCreationDTO.getName());

        // Create and save the task
        Task newTask = Task.builder()
                           .name(taskCreationDTO.getName())
                           .description(taskCreationDTO.getDescription())
                           .user(user)
                           .build();

        Task savedTask = taskRepository.save(newTask);
        log.info("Task '{}' created successfully for user with ID: {}", savedTask.getName(), user.getUserId());

        return TaskDto.builder()
                      .taskId(savedTask.getTaskId())
                      .name(savedTask.getName())
                      .description(savedTask.getDescription())
                      .userId(user.getUserId())
                      .build();
    }

    /**
     * Updates an existing task for a user, ensuring the task belongs to the user.
     *
     * @param userId        the ID of the user updating the task
     * @param taskId        the ID of the task to update
     * @param taskUpdateDTO the DTO containing updated task details
     *
     * @return the updated {@link TaskDto}
     *
     * @throws TaskNotFoundException    if the task does not exist
     * @throws IllegalArgumentException if the task does not belong to the user
     */
    public TaskDto updateTask(Long userId, Long taskId, TaskUpdateDTO taskUpdateDTO) {
        log.info("Updating task with ID: {} for user ID: {}", taskId, userId);
        TaskUtils.validateTaskUpdateDTO(taskUpdateDTO);

        // Fetch task by ID
        Task task = taskRepository.findById(taskId).orElseThrow(() -> new TaskNotFoundException(TASK_NOT_FOUND_EX));

        // Verify that the task belongs to the user
        if (!task.getUser().getUserId().equals(userId)) {
            log.warn("Update. Task with ID: {} does not belong to user with ID: {}", taskId, userId);
            throw new AccessDeniedException("Task does not belong to the user");
        }

        // Update task details
        task.setName(taskUpdateDTO.getName());
        task.setDescription(taskUpdateDTO.getDescription());

        // Save updated task
        Task updatedTask = taskRepository.save(task);
        log.info("Task with ID: {} updated successfully", updatedTask.getTaskId());

        return TaskDto.builder()
                      .taskId(updatedTask.getTaskId())
                      .name(updatedTask.getName())
                      .description(updatedTask.getDescription())
                      .userId(updatedTask.getUser().getUserId())
                      .build();
    }

    /**
     * Deletes a task if it belongs to the specified user.
     *
     * @param userId the ID of the user attempting to delete the task
     * @param taskId the ID of the task to delete
     *
     * @throws TaskNotFoundException if the task does not exist
     * @throws AccessDeniedException if the task does not belong to the user
     */
    public void deleteTask(Long userId, Long taskId) {
        log.info("Deleting task with ID: {} for user ID: {}", taskId, userId);

        Task task = taskRepository.findById(taskId).orElseThrow(() -> new TaskNotFoundException(TASK_NOT_FOUND_EX));

        if (!task.getUser().getUserId().equals(userId)) {
            log.warn("Delete. Task with ID: {} does not belong to user with ID: {}", taskId, userId);
            throw new AccessDeniedException("Task does not belong to the user");
        }

        taskRepository.delete(task);
        log.info("Task with ID: {} deleted successfully for user ID: {}", taskId, userId);
    }

    /**
     * Retrieves all tasks for a specific user.
     *
     * @param userId the ID of the user whose tasks are being retrieved
     *
     * @return a list of {@link TaskDto} representing all tasks for the user
     *
     * @throws IllegalArgumentException if the user does not exist
     */
    public List<TaskDto> getAllTasksForUser(Long userId) {
        log.debug("Fetching all tasks for user with ID: {}", userId);

        User user = userRepository.findById(userId).orElseThrow(() -> new IllegalArgumentException(USER_NOT_FOUND_EX));

        List<TaskDto> tasks = taskRepository.findAllByUser(user)
                                            .stream()
                                            .map(task -> TaskDto.builder()
                                                                .taskId(task.getTaskId())
                                                                .name(task.getName())
                                                                .description(task.getDescription())
                                                                .userId(task.getUser().getUserId())
                                                                .build())
                                            .toList();

        log.info("Found {} tasks for user with ID: {}", tasks.size(), userId);
        return tasks;
    }

    /**
     * Retrieves a specific task by ID for a user, ensuring the task belongs to the user.
     *
     * @param userId the ID of the user
     * @param taskId the ID of the task
     *
     * @return an {@link Optional} containing the task data if found, or empty if not found
     *
     * @throws IllegalArgumentException if the user does not exist
     */
    public Optional<TaskDto> getTaskByUserIdAndTaskId(Long userId, Long taskId) {
        log.debug("Searching for task with ID: {} for user with ID: {}", taskId, userId);

        User user = userRepository.findById(userId).orElseThrow(() -> new IllegalArgumentException(USER_NOT_FOUND_EX));

        return taskRepository.findByUserAndTaskId(user, taskId).map(task -> {
            log.info("Task found with ID: {} for user with ID: {}", taskId, userId);
            return TaskDto.builder()
                          .taskId(task.getTaskId())
                          .name(task.getName())
                          .description(task.getDescription())
                          .userId(user.getUserId())
                          .build();
        });
    }
}
