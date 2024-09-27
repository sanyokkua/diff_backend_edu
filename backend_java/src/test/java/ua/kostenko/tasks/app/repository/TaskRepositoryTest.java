package ua.kostenko.tasks.app.repository;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.jdbc.AutoConfigureTestDatabase;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;
import org.springframework.context.annotation.Import;
import ua.kostenko.tasks.app.TestcontainersConfiguration;
import ua.kostenko.tasks.app.entity.Task;
import ua.kostenko.tasks.app.entity.User;

import java.util.List;
import java.util.Optional;

import static org.assertj.core.api.AssertionsForClassTypes.assertThat;
import static org.junit.jupiter.api.Assertions.*;

@Import(TestcontainersConfiguration.class)
@DataJpaTest
@AutoConfigureTestDatabase(replace = AutoConfigureTestDatabase.Replace.NONE)
class TaskRepositoryTest {

    @Autowired
    private TaskRepository taskRepository;
    @Autowired
    private UserRepository userRepository;

    private User createAndSaveUser() {
        User user = User.builder().email("testuser@example.com").passwordHash("password123").build();
        return userRepository.save(user);
    }

    private Task createAndSaveTask(User user, String name, String description) {
        Task task = Task.builder().name(name).description(description).user(user).build();
        return taskRepository.save(task);
    }

    @Test
    void testSaveTask_Success() {
        User user = createAndSaveUser();
        Task task = Task.builder().name("Sample Task").description("Task Description").user(user).build();
        Task savedTask = taskRepository.save(task);

        assertThat(savedTask.getTaskId()).isNotNull();
        assertThat(savedTask.getName()).isEqualTo("Sample Task");
        assertThat(savedTask.getUser().getUserId()).isEqualTo(user.getUserId());
    }

    @Test
    void testFindById_Success() {
        User user = createAndSaveUser();
        Task task = createAndSaveTask(user, "Sample Task", "Task Description");
        Optional<Task> foundTask = taskRepository.findById(task.getTaskId());

        assertTrue(foundTask.isPresent());
        assertThat(foundTask.get().getName()).isEqualTo("Sample Task");
    }

    @Test
    void testFindById_TaskNotFound() {
        Optional<Task> foundTask = taskRepository.findById(999L);

        assertFalse(foundTask.isPresent());
    }

    @Test
    void testUpdateTask_Success() {
        User user = createAndSaveUser();
        Task task = createAndSaveTask(user, "Original Task", "Original Description");
        task.setName("Updated Task");
        Task updatedTask = taskRepository.save(task);

        assertThat(updatedTask.getTaskId()).isEqualTo(task.getTaskId());
        assertThat(updatedTask.getName()).isEqualTo("Updated Task");
    }

    @Test
    void testDeleteTask_Success() {
        User user = createAndSaveUser();
        Task task = createAndSaveTask(user, "Task to be deleted", "Description");
        taskRepository.delete(task);
        Optional<Task> deletedTask = taskRepository.findById(task.getTaskId());

        assertFalse(deletedTask.isPresent());
    }

    @Test
    void testFindByUserAndName_Success() {
        User user = createAndSaveUser();
        createAndSaveTask(user, "Unique Task", "Description");

        Optional<Task> foundTask = taskRepository.findByUserAndName(user, "Unique Task");

        assertTrue(foundTask.isPresent());
        assertThat(foundTask.get().getName()).isEqualTo("Unique Task");
    }

    @Test
    void testFindByUserAndName_TaskNotFound() {
        User user = createAndSaveUser();
        Optional<Task> foundTask = taskRepository.findByUserAndName(user, "Nonexistent Task");

        assertFalse(foundTask.isPresent());
    }

    @Test
    void testFindByUserAndTaskId_Success() {
        User user = createAndSaveUser();
        Task task = createAndSaveTask(user, "Sample Task", "Description");
        Optional<Task> foundTask = taskRepository.findByUserAndTaskId(user, task.getTaskId());

        assertTrue(foundTask.isPresent());
        assertThat(foundTask.get().getTaskId()).isEqualTo(task.getTaskId());
    }

    @Test
    void testFindByUserAndTaskId_TaskNotFound() {
        User user = createAndSaveUser();
        Optional<Task> foundTask = taskRepository.findByUserAndTaskId(user, 999L);

        assertFalse(foundTask.isPresent());
    }

    @Test
    void testFindAllByUser_Success() {
        User user = createAndSaveUser();
        createAndSaveTask(user, "Task 1", "Description 1");
        createAndSaveTask(user, "Task 2", "Description 2");
        List<Task> userTasks = taskRepository.findAllByUser(user);

        assertEquals(2, userTasks.size());
    }

    @Test
    void testFindAllByUser_NoTasks() {
        User user = createAndSaveUser();
        List<Task> userTasks = taskRepository.findAllByUser(user);

        assertTrue(userTasks.isEmpty());
    }

    @Test
    void testSaveTask_NullName_ThrowsException() {
        User user = createAndSaveUser();
        Task task = Task.builder().name(null).description("Task Description").user(user).build();

        assertThrows(Exception.class, () -> taskRepository.save(task));
    }

    @Test
    void testSaveTask_NullUser_ThrowsException() {
        Task task = Task.builder().name("Task Name").description("Task Description").user(null).build();

        assertThrows(Exception.class, () -> taskRepository.save(task));
    }
}