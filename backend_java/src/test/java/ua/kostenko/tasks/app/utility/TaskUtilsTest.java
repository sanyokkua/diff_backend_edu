package ua.kostenko.tasks.app.utility;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import ua.kostenko.tasks.app.dto.task.TaskCreationDTO;
import ua.kostenko.tasks.app.dto.task.TaskUpdateDTO;
import ua.kostenko.tasks.app.entity.Task;
import ua.kostenko.tasks.app.entity.User;
import ua.kostenko.tasks.app.exception.TaskAlreadyExistsException;
import ua.kostenko.tasks.app.repository.TaskRepository;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class TaskUtilsTest {

    @Mock
    private TaskRepository taskRepository;

    private User user;
    private Task task;

    @BeforeEach
    void setUp() {
        // Initialize test user and task
        user = new User();
        user.setUserId(1L);

        task = new Task();
        task.setName("Test Task");
    }

    // Test cases for checkTaskExistsForUser method
    @Test
    void checkTaskExistsForUser_taskAlreadyExists_shouldThrowException() {
        // Arrange
        when(taskRepository.findByUserAndName(user, "Test Task")).thenReturn(Optional.of(task));

        // Act & Assert
        TaskAlreadyExistsException exception = assertThrows(TaskAlreadyExistsException.class,
                                                            () -> TaskUtils.checkTaskExistsForUser(taskRepository,
                                                                                                   user,
                                                                                                   "Test Task"));
        assertEquals("Task with the name 'Test Task' already exists for the user", exception.getMessage());

        verify(taskRepository).findByUserAndName(user, "Test Task");
    }

    @Test
    void checkTaskExistsForUser_taskDoesNotExist_shouldNotThrowException() {
        // Arrange
        when(taskRepository.findByUserAndName(user, "Test Task")).thenReturn(Optional.empty());

        // Act & Assert
        assertDoesNotThrow(() -> TaskUtils.checkTaskExistsForUser(taskRepository, user, "Test Task"));

        verify(taskRepository).findByUserAndName(user, "Test Task");
    }

    @Test
    void checkTaskExistsForUser_nullUser_shouldThrowException() {
        // Act & Assert
        assertThrows(NullPointerException.class,
                     () -> TaskUtils.checkTaskExistsForUser(taskRepository, null, "Test Task"));
    }

    @Test
    void checkTaskExistsForUser_nullTaskName_shouldNotThrowException() {
        // Arrange
        when(taskRepository.findByUserAndName(user, null)).thenReturn(Optional.empty());

        // Act & Assert
        assertDoesNotThrow(() -> TaskUtils.checkTaskExistsForUser(taskRepository, user, null));

        verify(taskRepository).findByUserAndName(user, null);
    }

    // Test cases for validateTaskCreation method
    @Test
    void validateTaskCreation_validDTO_shouldNotThrowException() {
        // Arrange
        TaskCreationDTO taskCreationDTO = new TaskCreationDTO("Valid Task", "Valid Description");

        // Act & Assert
        assertDoesNotThrow(() -> TaskUtils.validateTaskCreation(taskCreationDTO));
    }

    @Test
    void validateTaskCreation_nullDTO_shouldThrowException() {
        // Act & Assert
        IllegalArgumentException exception =
                assertThrows(IllegalArgumentException.class, () -> TaskUtils.validateTaskCreation(null));
        assertEquals("TaskCreationDTO is null", exception.getMessage());
    }

    @Test
    void validateTaskCreation_nullName_shouldThrowException() {
        // Arrange
        TaskCreationDTO taskCreationDTO = new TaskCreationDTO(null, "Valid Description");

        // Act & Assert
        IllegalArgumentException exception =
                assertThrows(IllegalArgumentException.class, () -> TaskUtils.validateTaskCreation(taskCreationDTO));
        assertEquals("Task name cannot be null or empty", exception.getMessage());
    }

    @Test
    void validateTaskCreation_emptyName_shouldThrowException() {
        // Arrange
        TaskCreationDTO taskCreationDTO = new TaskCreationDTO("  ", "Valid Description");

        // Act & Assert
        IllegalArgumentException exception =
                assertThrows(IllegalArgumentException.class, () -> TaskUtils.validateTaskCreation(taskCreationDTO));
        assertEquals("Task name cannot be null or empty", exception.getMessage());
    }

    @Test
    void validateTaskCreation_nullDescription_shouldThrowException() {
        // Arrange
        TaskCreationDTO taskCreationDTO = new TaskCreationDTO("Valid Task", null);

        // Act & Assert
        IllegalArgumentException exception =
                assertThrows(IllegalArgumentException.class, () -> TaskUtils.validateTaskCreation(taskCreationDTO));
        assertEquals("Task description cannot be null or empty", exception.getMessage());
    }

    @Test
    void validateTaskCreation_emptyDescription_shouldThrowException() {
        // Arrange
        TaskCreationDTO taskCreationDTO = new TaskCreationDTO("Valid Task", "  ");

        // Act & Assert
        IllegalArgumentException exception =
                assertThrows(IllegalArgumentException.class, () -> TaskUtils.validateTaskCreation(taskCreationDTO));
        assertEquals("Task description cannot be null or empty", exception.getMessage());
    }

    // Test cases for validateTaskUpdateDTO method
    @Test
    void validateTaskUpdateDTO_validDTO_shouldNotThrowException() {
        // Arrange
        TaskUpdateDTO taskUpdateDTO = new TaskUpdateDTO("Valid Task", "Valid Description");

        // Act & Assert
        assertDoesNotThrow(() -> TaskUtils.validateTaskUpdateDTO(taskUpdateDTO));
    }

    @Test
    void validateTaskUpdateDTO_nullDTO_shouldThrowException() {
        // Act & Assert
        IllegalArgumentException exception =
                assertThrows(IllegalArgumentException.class, () -> TaskUtils.validateTaskUpdateDTO(null));
        assertEquals("TaskUpdateDTO is null", exception.getMessage());
    }

    @Test
    void validateTaskUpdateDTO_nullName_shouldThrowException() {
        // Arrange
        TaskUpdateDTO taskUpdateDTO = new TaskUpdateDTO(null, "Valid Description");

        // Act & Assert
        IllegalArgumentException exception =
                assertThrows(IllegalArgumentException.class, () -> TaskUtils.validateTaskUpdateDTO(taskUpdateDTO));
        assertEquals("Task name cannot be null or empty", exception.getMessage());
    }

    @Test
    void validateTaskUpdateDTO_emptyName_shouldThrowException() {
        // Arrange
        TaskUpdateDTO taskUpdateDTO = new TaskUpdateDTO("  ", "Valid Description");

        // Act & Assert
        IllegalArgumentException exception =
                assertThrows(IllegalArgumentException.class, () -> TaskUtils.validateTaskUpdateDTO(taskUpdateDTO));
        assertEquals("Task name cannot be null or empty", exception.getMessage());
    }

    @Test
    void validateTaskUpdateDTO_nullDescription_shouldThrowException() {
        // Arrange
        TaskUpdateDTO taskUpdateDTO = new TaskUpdateDTO("Valid Task", null);

        // Act & Assert
        IllegalArgumentException exception =
                assertThrows(IllegalArgumentException.class, () -> TaskUtils.validateTaskUpdateDTO(taskUpdateDTO));
        assertEquals("Task description cannot be null", exception.getMessage());
    }
}