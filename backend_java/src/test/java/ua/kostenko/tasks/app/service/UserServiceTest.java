package ua.kostenko.tasks.app.service;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.security.crypto.password.PasswordEncoder;
import ua.kostenko.tasks.app.dto.user.UserCreationDTO;
import ua.kostenko.tasks.app.dto.user.UserDeletionDTO;
import ua.kostenko.tasks.app.dto.user.UserDto;
import ua.kostenko.tasks.app.dto.user.UserUpdateDTO;
import ua.kostenko.tasks.app.entity.User;
import ua.kostenko.tasks.app.exception.EmailAlreadyExistsException;
import ua.kostenko.tasks.app.exception.InvalidEmailFormatException;
import ua.kostenko.tasks.app.exception.InvalidPasswordException;
import ua.kostenko.tasks.app.repository.UserRepository;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class UserServiceTest {

    @Mock
    private UserRepository userRepository;

    @Mock
    private PasswordEncoder passwordEncoder;

    @InjectMocks
    private UserService userService;

    private UserCreationDTO validUserCreationDTO;
    private UserUpdateDTO validUserUpdateDTO;
    private UserDeletionDTO validUserDeletionDTO;
    private User mockUser;

    @BeforeEach
    void setup() {
        validUserCreationDTO = UserCreationDTO.builder()
                                              .email("test@example.com")
                                              .password("password123")
                                              .passwordConfirmation("password123")
                                              .build();

        validUserUpdateDTO = UserUpdateDTO.builder()
                                          .currentPassword("oldPassword123")
                                          .newPassword("newPassword123")
                                          .newPasswordConfirmation("newPassword123")
                                          .build();

        validUserDeletionDTO =
                UserDeletionDTO.builder().email("test@example.com").currentPassword("password123").build();

        mockUser = User.builder().userId(1L).email("test@example.com").passwordHash("hashedPassword").build();
    }

    // Test for createUser
    @Test
    void testCreateUser_Success() {
        // Mock repository interactions
        when(userRepository.findByEmail(validUserCreationDTO.getEmail())).thenReturn(Optional.empty());
        when(passwordEncoder.encode(validUserCreationDTO.getPassword())).thenReturn("hashedPassword");
        when(userRepository.save(any(User.class))).thenReturn(mockUser);

        UserDto createdUser = userService.createUser(validUserCreationDTO);

        // Verify repository interactions
        verify(userRepository).save(any(User.class));
        assertEquals("test@example.com", createdUser.getEmail());
        assertEquals(1L, createdUser.getUserId());
    }

    @Test
    void testCreateUser_EmailAlreadyExists_ThrowsException() {
        // Simulate existing email in repository
        when(userRepository.findByEmail(validUserCreationDTO.getEmail())).thenReturn(Optional.of(mockUser));

        assertThrows(EmailAlreadyExistsException.class, () -> userService.createUser(validUserCreationDTO));

        verify(userRepository, never()).save(any(User.class));
    }

    @Test
    void testCreateUser_InvalidEmailFormat_ThrowsException() {
        validUserCreationDTO.setEmail("invalid-email");

        assertThrows(InvalidEmailFormatException.class, () -> userService.createUser(validUserCreationDTO));

        verify(userRepository, never()).save(any(User.class));
    }

    @Test
    void testCreateUser_InvalidPasswords_ThrowsException() {
        validUserCreationDTO.setPasswordConfirmation("wrongPassword");

        assertThrows(InvalidPasswordException.class, () -> userService.createUser(validUserCreationDTO));

        verify(userRepository, never()).save(any(User.class));
    }

    // Test for updateUserPassword
    @Test
    void testUpdateUserPassword_Success() {
        when(userRepository.findById(1L)).thenReturn(Optional.of(mockUser));
        when(passwordEncoder.matches(validUserUpdateDTO.getCurrentPassword(), mockUser.getPasswordHash())).thenReturn(
                true);
        when(passwordEncoder.encode(validUserUpdateDTO.getNewPassword())).thenReturn("newHashedPassword");
        when(userRepository.save(mockUser)).thenReturn(mockUser);

        UserDto updatedUser = userService.updateUserPassword(1L, validUserUpdateDTO);

        assertEquals(1L, updatedUser.getUserId());
        assertEquals("test@example.com", updatedUser.getEmail());

        verify(userRepository).save(mockUser);
    }

    @Test
    void testUpdateUserPassword_UserNotFound_ThrowsException() {
        when(userRepository.findById(1L)).thenReturn(Optional.empty());

        assertThrows(IllegalArgumentException.class, () -> userService.updateUserPassword(1L, validUserUpdateDTO));

        verify(userRepository, never()).save(any(User.class));
    }

    @Test
    void testUpdateUserPassword_InvalidCurrentPassword_ThrowsException() {
        when(userRepository.findById(1L)).thenReturn(Optional.of(mockUser));
        when(passwordEncoder.matches(validUserUpdateDTO.getCurrentPassword(), mockUser.getPasswordHash())).thenReturn(
                false);

        assertThrows(InvalidPasswordException.class, () -> userService.updateUserPassword(1L, validUserUpdateDTO));

        verify(userRepository, never()).save(any(User.class));
    }

    // Test for deleteUser
    @Test
    void testDeleteUser_Success() {
        when(userRepository.findById(1L)).thenReturn(Optional.of(mockUser));
        when(passwordEncoder.matches(validUserDeletionDTO.getCurrentPassword(), mockUser.getPasswordHash())).thenReturn(
                true);

        userService.deleteUser(1L, validUserDeletionDTO);

        verify(userRepository).delete(mockUser);
    }

    @Test
    void testDeleteUser_UserNotFound_ThrowsException() {
        when(userRepository.findById(1L)).thenReturn(Optional.empty());

        assertThrows(IllegalArgumentException.class, () -> userService.deleteUser(1L, validUserDeletionDTO));

        verify(userRepository, never()).delete(any(User.class));
    }

    @Test
    void testDeleteUser_InvalidPassword_ThrowsException() {
        when(userRepository.findById(1L)).thenReturn(Optional.of(mockUser));
        when(passwordEncoder.matches(validUserDeletionDTO.getCurrentPassword(), mockUser.getPasswordHash())).thenReturn(
                false);

        assertThrows(InvalidPasswordException.class, () -> userService.deleteUser(1L, validUserDeletionDTO));

        verify(userRepository, never()).delete(any(User.class));
    }
}