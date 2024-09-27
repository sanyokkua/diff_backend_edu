package ua.kostenko.tasks.app.repository;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.jdbc.AutoConfigureTestDatabase;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;
import org.springframework.context.annotation.Import;
import ua.kostenko.tasks.app.TestcontainersConfiguration;
import ua.kostenko.tasks.app.entity.User;

import java.util.Optional;

import static org.assertj.core.api.AssertionsForClassTypes.assertThat;
import static org.junit.jupiter.api.Assertions.*;

@Import(TestcontainersConfiguration.class)
@DataJpaTest
@AutoConfigureTestDatabase(replace = AutoConfigureTestDatabase.Replace.NONE)
class UserRepositoryTest {

    @Autowired
    private UserRepository userRepository;

    @Test
    void testSaveUser_Success() {
        User user = User.builder().email("testuser@example.com").passwordHash("password123").build();

        User savedUser = userRepository.save(user);

        assertThat(savedUser.getUserId()).isNotNull();
        assertThat(savedUser.getEmail()).isEqualTo("testuser@example.com");
    }

    @Test
    void testFindUserById_Success() {
        User user = User.builder().email("testuser@example.com").passwordHash("password123").build();
        User savedUser = userRepository.save(user);

        Optional<User> foundUser = userRepository.findById(savedUser.getUserId());

        assertTrue(foundUser.isPresent());
        assertThat(foundUser.get().getEmail()).isEqualTo("testuser@example.com");
    }

    @Test
    void testUpdateUser_Success() {
        User user = User.builder().email("testuser@example.com").passwordHash("password123").build();
        User savedUser = userRepository.save(user);

        savedUser.setEmail("updateduser@example.com");
        User updatedUser = userRepository.save(savedUser);

        assertThat(updatedUser.getUserId()).isEqualTo(savedUser.getUserId());
        assertThat(updatedUser.getEmail()).isEqualTo("updateduser@example.com");
    }

    @Test
    void testDeleteUser_Success() {
        User user = User.builder().email("testuser@example.com").passwordHash("password123").build();
        User savedUser = userRepository.save(user);

        userRepository.delete(savedUser);
        Optional<User> deletedUser = userRepository.findById(savedUser.getUserId());

        assertFalse(deletedUser.isPresent());
    }

    @Test
    void testFindByEmail_Success() {
        User user = User.builder().email("testuser@example.com").passwordHash("password123").build();
        userRepository.save(user);

        Optional<User> foundUser = userRepository.findByEmail("testuser@example.com");

        assertTrue(foundUser.isPresent());
        assertThat(foundUser.get().getEmail()).isEqualTo("testuser@example.com");
    }

    @Test
    void testFindByEmail_UserNotFound() {
        Optional<User> foundUser = userRepository.findByEmail("nonexistent@example.com");

        assertFalse(foundUser.isPresent());
    }

    @Test
    void testFindUserById_UserNotFound() {
        Optional<User> foundUser = userRepository.findById(999L);

        assertFalse(foundUser.isPresent());
    }

    @Test
    void testSaveUser_NullEmail_ThrowsException() {
        User user = User.builder().email(null).passwordHash("password123").build();

        assertThrows(Exception.class, () -> userRepository.save(user));
    }

    @Test
    void testSaveUser_DuplicateEmail_ThrowsException() {
        User user1 = User.builder().email("duplicate@example.com").passwordHash("password123").build();
        User user2 = User.builder().email("duplicate@example.com").passwordHash("password456").build();

        userRepository.save(user1);

        assertThrows(Exception.class, () -> userRepository.save(user2));
    }
}