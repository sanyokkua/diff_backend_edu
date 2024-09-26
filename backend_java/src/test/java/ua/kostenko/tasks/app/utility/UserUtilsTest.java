package ua.kostenko.tasks.app.utility;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.security.access.AccessDeniedException;
import org.springframework.security.crypto.password.PasswordEncoder;
import ua.kostenko.tasks.app.dto.user.*;
import ua.kostenko.tasks.app.entity.User;
import ua.kostenko.tasks.app.exception.EmailAlreadyExistsException;
import ua.kostenko.tasks.app.exception.InvalidEmailFormatException;
import ua.kostenko.tasks.app.exception.InvalidPasswordException;
import ua.kostenko.tasks.app.repository.UserRepository;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class UserUtilsTest {

    private final UserRepository userRepository = mock(UserRepository.class);
    private final PasswordEncoder passwordEncoder = mock(PasswordEncoder.class);

    // Test for validateEmailFormat
    @Test
    void validateEmailFormat_validEmail() {
        assertDoesNotThrow(() -> UserUtils.validateEmailFormat("test@example.com"));
    }

    @Test
    void validateEmailFormat_invalidEmail() {
        assertThrows(InvalidEmailFormatException.class, () -> UserUtils.validateEmailFormat("invalid-email"));
    }

    @Test
    void validateEmailFormat_nullEmail() {
        assertThrows(InvalidEmailFormatException.class, () -> UserUtils.validateEmailFormat(null));
    }

    // Test for checkUserExists
    @Test
    void checkUserExists_userExists() {
        when(userRepository.findByEmail("test@example.com")).thenReturn(Optional.of(new User()));

        assertThrows(EmailAlreadyExistsException.class,
                     () -> UserUtils.checkUserExists(userRepository, "test@example.com"));

        verify(userRepository, times(1)).findByEmail("test@example.com");
    }

    @Test
    void checkUserExists_userDoesNotExist() {
        when(userRepository.findByEmail("test@example.com")).thenReturn(Optional.empty());

        assertDoesNotThrow(() -> UserUtils.checkUserExists(userRepository, "test@example.com"));

        verify(userRepository, times(1)).findByEmail("test@example.com");
    }

    // Test for validatePasswords
    @Test
    void validatePasswords_match() {
        assertDoesNotThrow(() -> UserUtils.validatePasswords("password", "password"));
    }

    @Test
    void validatePasswords_mismatch() {
        assertThrows(InvalidPasswordException.class, () -> UserUtils.validatePasswords("password", "different"));
    }

    // Test for validatePasswordUpdate
    @Test
    void validatePasswordUpdate_successful() {
        User user = new User();
        user.setEmail("test@example.com");
        user.setPasswordHash("hashedPassword");

        UserUpdateDTO userUpdateDTO = new UserUpdateDTO("currentPassword", "newPassword", "newPassword");

        when(passwordEncoder.matches("currentPassword", "hashedPassword")).thenReturn(true);

        assertDoesNotThrow(() -> UserUtils.validatePasswordUpdate(userUpdateDTO, user, passwordEncoder));

        verify(passwordEncoder, times(1)).matches("currentPassword", "hashedPassword");
    }

    @Test
    void validatePasswordUpdate_incorrectCurrentPassword() {
        User user = new User();
        user.setEmail("test@example.com");
        user.setPasswordHash("hashedPassword");

        UserUpdateDTO userUpdateDTO = new UserUpdateDTO("wrongPassword", "newPassword", "newPassword");

        when(passwordEncoder.matches("wrongPassword", "hashedPassword")).thenReturn(false);

        assertThrows(InvalidPasswordException.class,
                     () -> UserUtils.validatePasswordUpdate(userUpdateDTO, user, passwordEncoder));

        verify(passwordEncoder, times(1)).matches("wrongPassword", "hashedPassword");
    }

    @Test
    void validatePasswordUpdate_newPasswordSameAsCurrent() {
        User user = new User();
        user.setEmail("test@example.com");
        user.setPasswordHash("hashedPassword");

        UserUpdateDTO userUpdateDTO = new UserUpdateDTO("currentPassword", "currentPassword", "currentPassword");

        when(passwordEncoder.matches("currentPassword", "hashedPassword")).thenReturn(true);

        assertThrows(InvalidPasswordException.class,
                     () -> UserUtils.validatePasswordUpdate(userUpdateDTO, user, passwordEncoder));
    }

    @Test
    void validatePasswordUpdate_newPasswordMismatch() {
        User user = new User();
        user.setEmail("test@example.com");
        user.setPasswordHash("hashedPassword");

        UserUpdateDTO userUpdateDTO = new UserUpdateDTO("currentPassword", "newPassword", "differentNewPassword");

        when(passwordEncoder.matches("currentPassword", "hashedPassword")).thenReturn(true);

        assertThrows(InvalidPasswordException.class,
                     () -> UserUtils.validatePasswordUpdate(userUpdateDTO, user, passwordEncoder));
    }

    // Test for validateAuthenticatedUserIdWithPassed
    @Test
    void validateAuthenticatedUserIdWithPassed_successful() {
        UserDto userDto = new UserDto();
        userDto.setUserId(1L);

        assertDoesNotThrow(() -> UserUtils.validateAuthenticatedUserIdWithPassed(userDto, 1L));
    }

    @Test
    void validateAuthenticatedUserIdWithPassed_mismatch() {
        UserDto userDto = new UserDto();
        userDto.setUserId(1L);

        assertThrows(AccessDeniedException.class, () -> UserUtils.validateAuthenticatedUserIdWithPassed(userDto, 2L));
    }

    @Test
    void validateAuthenticatedUserIdWithPassed_nullUserDto() {
        assertThrows(AccessDeniedException.class, () -> UserUtils.validateAuthenticatedUserIdWithPassed(null, 1L));
    }

    @Test
    void validateAuthenticatedUserIdWithPassed_nullUserId() {
        UserDto userDto = new UserDto();
        userDto.setUserId(1L);

        assertThrows(AccessDeniedException.class, () -> UserUtils.validateAuthenticatedUserIdWithPassed(userDto, null));
    }

    // Test for validateUserLoginDto
    @Test
    void validateUserLoginDto_valid() {
        UserLoginDto userLoginDto = new UserLoginDto();
        userLoginDto.setEmail("test@example.com");
        userLoginDto.setPassword("password");

        assertDoesNotThrow(() -> UserUtils.validateUserLoginDto(userLoginDto));
    }

    @Test
    void validateUserLoginDto_nullDto() {
        assertThrows(IllegalArgumentException.class, () -> UserUtils.validateUserLoginDto(null));
    }

    @Test
    void validateUserLoginDto_nullEmail() {
        UserLoginDto userLoginDto = new UserLoginDto();
        userLoginDto.setPassword("password");

        assertThrows(IllegalArgumentException.class, () -> UserUtils.validateUserLoginDto(userLoginDto));
    }

    @Test
    void validateUserLoginDto_nullPassword() {
        UserLoginDto userLoginDto = new UserLoginDto();
        userLoginDto.setEmail("test@example.com");

        assertThrows(IllegalArgumentException.class, () -> UserUtils.validateUserLoginDto(userLoginDto));
    }

    // Test for validateUserCreationDTO
    @Test
    void validateUserCreationDTO_valid() {
        UserCreationDTO userCreationDTO = new UserCreationDTO();
        userCreationDTO.setEmail("test@example.com");
        userCreationDTO.setPassword("password");
        userCreationDTO.setPasswordConfirmation("password");

        assertDoesNotThrow(() -> UserUtils.validateUserCreationDTO(userCreationDTO));
    }

    @Test
    void validateUserCreationDTO_nullDto() {
        assertThrows(IllegalArgumentException.class, () -> UserUtils.validateUserCreationDTO(null));
    }

    @Test
    void validateUserCreationDTO_nullEmail() {
        UserCreationDTO userCreationDTO = new UserCreationDTO();
        userCreationDTO.setPassword("password");
        userCreationDTO.setPasswordConfirmation("password");

        assertThrows(IllegalArgumentException.class, () -> UserUtils.validateUserCreationDTO(userCreationDTO));
    }

    @Test
    void validateUserCreationDTO_nullPassword() {
        UserCreationDTO userCreationDTO = new UserCreationDTO();
        userCreationDTO.setEmail("test@example.com");
        userCreationDTO.setPasswordConfirmation("password");

        assertThrows(IllegalArgumentException.class, () -> UserUtils.validateUserCreationDTO(userCreationDTO));
    }

    @Test
    void validateUserCreationDTO_passwordMismatch() {
        UserCreationDTO userCreationDTO = new UserCreationDTO();
        userCreationDTO.setEmail("test@example.com");
        userCreationDTO.setPassword("password");
        userCreationDTO.setPasswordConfirmation("different");

        assertThrows(InvalidPasswordException.class, () -> UserUtils.validateUserCreationDTO(userCreationDTO));
    }

    @Test
    void validateUserCreationDTO_passwordConfirmationNull() {
        UserCreationDTO userCreationDTO = new UserCreationDTO();
        userCreationDTO.setEmail("test@example.com");
        userCreationDTO.setPassword("password");
        userCreationDTO.setPasswordConfirmation(null);

        assertThrows(IllegalArgumentException.class, () -> UserUtils.validateUserCreationDTO(userCreationDTO));
    }

    // Test for validateUserUpdateDTO
    @Test
    void validateUserUpdateDTO_valid() {
        UserUpdateDTO userUpdateDTO = new UserUpdateDTO("currentPassword", "newPassword", "newPassword");

        assertDoesNotThrow(() -> UserUtils.validateUserUpdateDTO(userUpdateDTO));
    }

    @Test
    void validateUserUpdateDTO_nullDto() {
        assertThrows(IllegalArgumentException.class, () -> UserUtils.validateUserUpdateDTO(null));
    }

    @Test
    void validateUserUpdateDTO_nullCurrentPassword() {
        UserUpdateDTO userUpdateDTO = new UserUpdateDTO(null, "newPassword", "newPassword");

        assertThrows(IllegalArgumentException.class, () -> UserUtils.validateUserUpdateDTO(userUpdateDTO));
    }

    @Test
    void validateUserUpdateDTO_nullNewPassword() {
        UserUpdateDTO userUpdateDTO = new UserUpdateDTO("currentPassword", null, "newPassword");

        assertThrows(IllegalArgumentException.class, () -> UserUtils.validateUserUpdateDTO(userUpdateDTO));
    }

    @Test
    void validateUserUpdateDTO_nullNewPasswordConfirm() {
        UserUpdateDTO userUpdateDTO = new UserUpdateDTO("currentPassword", "newPassword", null);

        assertThrows(IllegalArgumentException.class, () -> UserUtils.validateUserUpdateDTO(userUpdateDTO));
    }

    // Test for validateUserDeletionDTO
    @Test
    void validateUserDeletionDTO_valid() {
        UserDeletionDTO userDeletionDTO = new UserDeletionDTO("test@example.com", "password");

        assertDoesNotThrow(() -> UserUtils.validateUserDeletionDTO(userDeletionDTO));
    }

    @Test
    void validateUserDeletionDTO_nullDto() {
        assertThrows(IllegalArgumentException.class, () -> UserUtils.validateUserDeletionDTO(null));
    }

    @Test
    void validateUserDeletionDTO_nullEmail() {
        UserDeletionDTO userDeletionDTO = new UserDeletionDTO(null, "password");

        assertThrows(IllegalArgumentException.class, () -> UserUtils.validateUserDeletionDTO(userDeletionDTO));
    }

    @Test
    void validateUserDeletionDTO_nullPassword() {
        UserDeletionDTO userDeletionDTO = new UserDeletionDTO("test@example.com", null);

        assertThrows(IllegalArgumentException.class, () -> UserUtils.validateUserDeletionDTO(userDeletionDTO));
    }
}