package ua.kostenko.tasks.app.service;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.security.access.AccessDeniedException;
import ua.kostenko.tasks.app.dto.task.TaskCreationDTO;
import ua.kostenko.tasks.app.dto.task.TaskDto;
import ua.kostenko.tasks.app.dto.task.TaskUpdateDTO;
import ua.kostenko.tasks.app.entity.Task;
import ua.kostenko.tasks.app.entity.User;
import ua.kostenko.tasks.app.exception.TaskAlreadyExistsException;
import ua.kostenko.tasks.app.exception.TaskNotFoundException;
import ua.kostenko.tasks.app.repository.TaskRepository;
import ua.kostenko.tasks.app.repository.UserRepository;

import java.util.List;
import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class TaskServiceTest {

    @Mock
    private TaskRepository taskRepository;

    @Mock
    private UserRepository userRepository;

    @InjectMocks
    private TaskService taskService;

    private TaskCreationDTO validTaskCreationDTO;
    private TaskUpdateDTO validTaskUpdateDTO;
    private Task mockTask;
    private User mockUser;

    @BeforeEach
    void setup() {
        validTaskCreationDTO = TaskCreationDTO.builder().name("Sample Task").description("Sample Description").build();

        validTaskUpdateDTO = TaskUpdateDTO.builder().name("Updated Task").description("Updated Description").build();

        mockUser = User.builder().userId(1L).email("test@example.com").build();

        mockTask =
                Task.builder().taskId(1L).name("Sample Task").description("Sample Description").user(mockUser).build();
    }

    // Test for createTask
    @Test
    void testCreateTask_Success() {
        // Mock repository interactions
        when(userRepository.findById(1L)).thenReturn(Optional.of(mockUser));
        when(taskRepository.save(any(Task.class))).thenReturn(mockTask);

        // Actual call
        TaskDto createdTask = taskService.createTask(1L, validTaskCreationDTO);

        // Verify repository interactions
        verify(taskRepository).save(any(Task.class));
        assertEquals("Sample Task", createdTask.getName());
        assertEquals(1L, createdTask.getTaskId());
    }

    @Test
    void testCreateTask_UserNotFound_ThrowsException() {
        when(userRepository.findById(1L)).thenReturn(Optional.empty());

        assertThrows(IllegalArgumentException.class, () -> taskService.createTask(1L, validTaskCreationDTO));

        verify(taskRepository, never()).save(any(Task.class));
    }

    @Test
    void testCreateTask_TaskAlreadyExists_ThrowsException() {
        when(userRepository.findById(1L)).thenReturn(Optional.of(mockUser));
        doThrow(new TaskAlreadyExistsException("Task already exists")).when(taskRepository)
                                                                      .findByUserAndName(mockUser, "Sample Task");

        assertThrows(TaskAlreadyExistsException.class, () -> taskService.createTask(1L, validTaskCreationDTO));

        verify(taskRepository, never()).save(any(Task.class));
    }

    @Test
    void testCreateTask_InvalidTaskCreationDTO_ThrowsException() {
        TaskCreationDTO invalidDTO = TaskCreationDTO.builder().name("").description("").build();

        assertThrows(IllegalArgumentException.class, () -> taskService.createTask(1L, invalidDTO));

        verify(taskRepository, never()).save(any(Task.class));
    }

    // Test for updateTask
    @Test
    void testUpdateTask_Success() {
        when(taskRepository.findById(1L)).thenReturn(Optional.of(mockTask));
        when(taskRepository.save(any(Task.class))).thenReturn(mockTask);

        TaskDto updatedTask = taskService.updateTask(1L, 1L, validTaskUpdateDTO);

        verify(taskRepository).save(mockTask);
        assertEquals("Updated Task", updatedTask.getName());
    }

    @Test
    void testUpdateTask_TaskNotFound_ThrowsException() {
        when(taskRepository.findById(1L)).thenReturn(Optional.empty());

        assertThrows(TaskNotFoundException.class, () -> taskService.updateTask(1L, 1L, validTaskUpdateDTO));

        verify(taskRepository, never()).save(any(Task.class));
    }

    @Test
    void testUpdateTask_TaskNotBelongToUser_ThrowsAccessDeniedException() {
        User otherUser = User.builder().userId(2L).build();
        mockTask.setUser(otherUser);

        when(taskRepository.findById(1L)).thenReturn(Optional.of(mockTask));

        assertThrows(AccessDeniedException.class, () -> taskService.updateTask(1L, 1L, validTaskUpdateDTO));

        verify(taskRepository, never()).save(any(Task.class));
    }

    @Test
    void testUpdateTask_InvalidTaskUpdateDTO_ThrowsException() {
        TaskUpdateDTO invalidDTO = TaskUpdateDTO.builder().name("").description("").build();

        assertThrows(IllegalArgumentException.class, () -> taskService.updateTask(1L, 1L, invalidDTO));

        verify(taskRepository, never()).save(any(Task.class));
    }

    // Test for deleteTask
    @Test
    void testDeleteTask_Success() {
        when(taskRepository.findById(1L)).thenReturn(Optional.of(mockTask));

        taskService.deleteTask(1L, 1L);

        verify(taskRepository).delete(mockTask);
    }

    @Test
    void testDeleteTask_TaskNotFound_ThrowsException() {
        when(taskRepository.findById(1L)).thenReturn(Optional.empty());

        assertThrows(TaskNotFoundException.class, () -> taskService.deleteTask(1L, 1L));

        verify(taskRepository, never()).delete(any(Task.class));
    }

    @Test
    void testDeleteTask_TaskNotBelongToUser_ThrowsAccessDeniedException() {
        User otherUser = User.builder().userId(2L).build();
        mockTask.setUser(otherUser);

        when(taskRepository.findById(1L)).thenReturn(Optional.of(mockTask));

        assertThrows(AccessDeniedException.class, () -> taskService.deleteTask(1L, 1L));

        verify(taskRepository, never()).delete(any(Task.class));
    }

    // Test for getAllTasksForUser
    @Test
    void testGetAllTasksForUser_Success() {
        when(userRepository.findById(1L)).thenReturn(Optional.of(mockUser));
        when(taskRepository.findAllByUser(mockUser)).thenReturn(List.of(mockTask));

        List<TaskDto> tasks = taskService.getAllTasksForUser(1L);

        assertEquals(1, tasks.size());
        assertEquals("Sample Task", tasks.getFirst().getName());
    }

    @Test
    void testGetAllTasksForUser_UserNotFound_ThrowsException() {
        when(userRepository.findById(1L)).thenReturn(Optional.empty());

        assertThrows(IllegalArgumentException.class, () -> taskService.getAllTasksForUser(1L));

        verify(taskRepository, never()).findAllByUser(any(User.class));
    }

    // Test for getTaskByUserIdAndTaskId
    @Test
    void testGetTaskByUserIdAndTaskId_Success() {
        when(userRepository.findById(1L)).thenReturn(Optional.of(mockUser));
        when(taskRepository.findByUserAndTaskId(mockUser, 1L)).thenReturn(Optional.of(mockTask));

        Optional<TaskDto> task = taskService.getTaskByUserIdAndTaskId(1L, 1L);

        assertTrue(task.isPresent());
        assertEquals("Sample Task", task.get().getName());
    }

    @Test
    void testGetTaskByUserIdAndTaskId_UserNotFound_ThrowsException() {
        when(userRepository.findById(1L)).thenReturn(Optional.empty());

        assertThrows(IllegalArgumentException.class, () -> taskService.getTaskByUserIdAndTaskId(1L, 1L));

        verify(taskRepository, never()).findByUserAndTaskId(any(User.class), any(Long.class));
    }

    @Test
    void testGetTaskByUserIdAndTaskId_TaskNotFound_ReturnsEmpty() {
        when(userRepository.findById(1L)).thenReturn(Optional.of(mockUser));
        when(taskRepository.findByUserAndTaskId(mockUser, 1L)).thenReturn(Optional.empty());

        Optional<TaskDto> task = taskService.getTaskByUserIdAndTaskId(1L, 1L);

        assertFalse(task.isPresent());
    }
}