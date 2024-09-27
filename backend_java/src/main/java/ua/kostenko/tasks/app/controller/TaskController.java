package ua.kostenko.tasks.app.controller;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import ua.kostenko.tasks.app.dto.ResponseDto;
import ua.kostenko.tasks.app.dto.task.TaskCreationDTO;
import ua.kostenko.tasks.app.dto.task.TaskDto;
import ua.kostenko.tasks.app.dto.task.TaskUpdateDTO;
import ua.kostenko.tasks.app.dto.user.UserDto;
import ua.kostenko.tasks.app.service.AuthUserExtractionService;
import ua.kostenko.tasks.app.service.TaskService;
import ua.kostenko.tasks.app.utility.ResponseDtoUtils;
import ua.kostenko.tasks.app.utility.UserUtils;

import java.util.List;
import java.util.Optional;

/**
 * REST controller for managing tasks related to users.
 *
 * <p>This controller provides endpoints for authenticated users to create, retrieve, update,
 * and delete tasks associated with their account.</p>
 */
@Tag(name = "Task Management REST Controller", description = "Handles operations related to user tasks.")
@Slf4j
@RestController
@RequestMapping("api/v1/users/{userId}/tasks")
@RequiredArgsConstructor
public class TaskController {

    private final TaskService taskService;
    private final AuthUserExtractionService userExtractionService;

    /**
     * Create a new task for the authenticated user.
     *
     * @param userId          The ID of the user who owns the task.
     * @param taskCreationDTO The task creation data transfer object.
     *
     * @return ResponseEntity containing the created TaskDto or an error response.
     */
    @PostMapping("/")
    @Operation(summary = "Create New Task", description = "Creates a new task for the authenticated user.")
    public ResponseEntity<ResponseDto<TaskDto>> createTask(
            @Parameter(description = "The ID of the user who owns the task.") @PathVariable Long userId,
            @Parameter(description = "The task creation data transfer object.") @RequestBody TaskCreationDTO taskCreationDTO) {

        log.info("Request to create task for user with ID: {}", userId);

        // Extract and validate the authenticated user
        UserDto authenticatedUser = userExtractionService.getUserFromAuthContext();
        UserUtils.validateAuthenticatedUserIdWithPassed(authenticatedUser, userId);

        TaskDto createdTaskDto = taskService.createTask(userId, taskCreationDTO);
        log.info("Task '{}' created successfully for user ID: {}", taskCreationDTO.getName(), userId);
        return ResponseDtoUtils.buildDtoResponse(createdTaskDto, HttpStatus.CREATED);
    }

    /**
     * Retrieve a specific task by its ID.
     *
     * @param userId The ID of the user who owns the task.
     * @param taskId The ID of the task to retrieve.
     *
     * @return ResponseEntity containing the TaskDto or an error response.
     */
    @GetMapping("/{taskId}")
    @Operation(summary = "Retrieve Task by ID", description = "Retrieves a specific task by its ID for the authenticated user.")
    public ResponseEntity<ResponseDto<TaskDto>> getTaskById(
            @Parameter(description = "The ID of the user who owns the task.") @PathVariable Long userId,
            @Parameter(description = "The ID of the task to retrieve.") @PathVariable Long taskId) {

        log.info("Request to retrieve task ID: {} for user ID: {}", taskId, userId);

        // Extract and validate the authenticated user
        UserDto authenticatedUser = userExtractionService.getUserFromAuthContext();
        UserUtils.validateAuthenticatedUserIdWithPassed(authenticatedUser, userId);

        Optional<TaskDto> taskDto = taskService.getTaskByUserIdAndTaskId(userId, taskId);
        if (taskDto.isPresent()) {
            log.info("Task ID: {} found for user ID: {}", taskId, userId);
            return ResponseDtoUtils.buildDtoResponse(taskDto.get(), HttpStatus.OK);
        } else {
            log.warn("Task ID: {} not found for user ID: {}", taskId, userId);
            return ResponseDtoUtils.buildDtoResponse(null, HttpStatus.NOT_FOUND);
        }
    }

    /**
     * Update an existing task for the authenticated user.
     *
     * @param userId        The ID of the user who owns the task.
     * @param taskId        The ID of the task to update.
     * @param taskUpdateDTO The task update data transfer object.
     *
     * @return ResponseEntity containing the updated TaskDto or an error response.
     */
    @PutMapping("/{taskId}")
    @Operation(summary = "Update Task", description = "Updates an existing task for the authenticated user.")
    public ResponseEntity<ResponseDto<TaskDto>> updateTask(
            @Parameter(description = "The ID of the user who owns the task.") @PathVariable Long userId,
            @Parameter(description = "The ID of the task to update.") @PathVariable Long taskId,
            @Parameter(description = "The task update data transfer object.") @RequestBody TaskUpdateDTO taskUpdateDTO) {

        log.info("Request to update task ID: {} for user ID: {}", taskId, userId);

        // Extract and validate the authenticated user
        UserDto authenticatedUser = userExtractionService.getUserFromAuthContext();
        UserUtils.validateAuthenticatedUserIdWithPassed(authenticatedUser, userId);

        TaskDto updatedTaskDto = taskService.updateTask(userId, taskId, taskUpdateDTO);
        log.info("Task ID: {} updated successfully for user ID: {}", taskId, userId);
        return ResponseDtoUtils.buildDtoResponse(updatedTaskDto, HttpStatus.OK);
    }

    /**
     * Delete a specific task by its ID.
     *
     * @param userId The ID of the user who owns the task.
     * @param taskId The ID of the task to delete.
     *
     * @return ResponseEntity with HTTP status.
     */
    @DeleteMapping("/{taskId}")
    @Operation(summary = "Delete Task", description = "Deletes a specific task by its ID for the authenticated user.")
    public ResponseEntity<ResponseDto<Void>> deleteTask(
            @Parameter(description = "The ID of the user who owns the task.") @PathVariable Long userId,
            @Parameter(description = "The ID of the task to delete.") @PathVariable Long taskId) {

        log.info("Request to delete task ID: {} for user ID: {}", taskId, userId);

        // Extract and validate the authenticated user
        UserDto authenticatedUser = userExtractionService.getUserFromAuthContext();
        UserUtils.validateAuthenticatedUserIdWithPassed(authenticatedUser, userId);

        taskService.deleteTask(userId, taskId);
        log.info("Task ID: {} deleted successfully for user ID: {}", taskId, userId);
        return ResponseDtoUtils.buildDtoResponse(null, HttpStatus.NO_CONTENT);
    }

    /**
     * Retrieve all tasks for the authenticated user.
     *
     * @param userId The ID of the user who owns the tasks.
     *
     * @return ResponseEntity containing a list of TaskDto or an error response.
     */
    @GetMapping("/")
    @Operation(summary = "Retrieve All Tasks", description = "Retrieves all tasks for the authenticated user.")
    public ResponseEntity<ResponseDto<List<TaskDto>>> getAllTasksForUser(
            @Parameter(description = "The ID of the user who owns the tasks.") @PathVariable Long userId) {

        log.info("Request to retrieve all tasks for user ID: {}", userId);

        // Extract and validate the authenticated user
        UserDto authenticatedUser = userExtractionService.getUserFromAuthContext();
        UserUtils.validateAuthenticatedUserIdWithPassed(authenticatedUser, userId);

        List<TaskDto> tasks = taskService.getAllTasksForUser(userId);
        log.info("{} tasks retrieved for user ID: {}", tasks.size(), userId);
        return ResponseDtoUtils.buildDtoResponse(tasks, HttpStatus.OK);
    }
}
