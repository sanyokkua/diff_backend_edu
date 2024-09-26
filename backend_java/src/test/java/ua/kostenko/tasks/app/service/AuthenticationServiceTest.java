package ua.kostenko.tasks.app.service;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.security.crypto.password.PasswordEncoder;
import ua.kostenko.tasks.app.dto.user.UserCreationDTO;
import ua.kostenko.tasks.app.dto.user.UserDto;
import ua.kostenko.tasks.app.dto.user.UserLoginDto;
import ua.kostenko.tasks.app.entity.User;
import ua.kostenko.tasks.app.exception.InvalidPasswordException;
import ua.kostenko.tasks.app.repository.UserRepository;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class AuthenticationServiceTest {

    private static final String VALID_EMAIL = "test@example.com";
    private static final String INVALID_EMAIL = "invalid@example.com";
    private static final String VALID_PASSWORD = "validPassword";
    private static final String INVALID_PASSWORD = "invalidPassword";
    private static final String ENCODED_PASSWORD = "encodedPassword";
    private static final String JWT_TOKEN = "jwtToken";
    @Mock
    private UserService userService;
    @Mock
    private UserRepository userRepository;
    @Mock
    private JwtService jwtService;
    @Mock
    private PasswordEncoder passwordEncoder;
    @InjectMocks
    private AuthenticationService authenticationService;
    private User user;

    @BeforeEach
    void setup() {
        user = User.builder().email(VALID_EMAIL).passwordHash(ENCODED_PASSWORD).userId(1L).build();
    }

    // Tests for loginUser method

    @Test
    void loginUser_successfulLogin_returnsUserDtoWithJwtToken() {
        UserLoginDto userLoginDto = new UserLoginDto(VALID_EMAIL, VALID_PASSWORD);

        // Mocking repository, passwordEncoder, and jwtService behaviors
        when(userRepository.findByEmail(VALID_EMAIL)).thenReturn(Optional.of(user));
        when(passwordEncoder.matches(VALID_PASSWORD, ENCODED_PASSWORD)).thenReturn(true);
        when(jwtService.generateJwtToken(VALID_EMAIL)).thenReturn(JWT_TOKEN);

        // Call the method to test
        UserDto result = authenticationService.loginUser(userLoginDto);

        // Assertions
        assertNotNull(result);
        assertEquals(user.getUserId(), result.getUserId());
        assertEquals(VALID_EMAIL, result.getEmail());
        assertEquals(JWT_TOKEN, result.getJwtToken());

        // Verify interactions
        verify(userRepository).findByEmail(VALID_EMAIL);
        verify(passwordEncoder).matches(VALID_PASSWORD, ENCODED_PASSWORD);
        verify(jwtService).generateJwtToken(VALID_EMAIL);
    }

    @Test
    void loginUser_userNotFound_throwsIllegalArgumentException() {
        UserLoginDto userLoginDto = new UserLoginDto(INVALID_EMAIL, VALID_PASSWORD);

        // Mock repository to return empty
        when(userRepository.findByEmail(INVALID_EMAIL)).thenReturn(Optional.empty());

        // Expect exception
        IllegalArgumentException exception =
                assertThrows(IllegalArgumentException.class, () -> authenticationService.loginUser(userLoginDto));

        assertEquals("User not found", exception.getMessage());

        // Verify repository was called
        verify(userRepository).findByEmail(INVALID_EMAIL);
        verify(passwordEncoder, never()).matches(anyString(), anyString());
        verify(jwtService, never()).generateJwtToken(anyString());
    }

    @Test
    void loginUser_invalidPassword_throwsInvalidPasswordException() {
        UserLoginDto userLoginDto = new UserLoginDto(VALID_EMAIL, INVALID_PASSWORD);

        // Mock repository and passwordEncoder behavior
        when(userRepository.findByEmail(VALID_EMAIL)).thenReturn(Optional.of(user));
        when(passwordEncoder.matches(INVALID_PASSWORD, ENCODED_PASSWORD)).thenReturn(false);

        // Expect exception
        InvalidPasswordException exception =
                assertThrows(InvalidPasswordException.class, () -> authenticationService.loginUser(userLoginDto));

        assertEquals("Invalid credentials", exception.getMessage());

        // Verify interactions
        verify(userRepository).findByEmail(VALID_EMAIL);
        verify(passwordEncoder).matches(INVALID_PASSWORD, ENCODED_PASSWORD);
        verify(jwtService, never()).generateJwtToken(anyString());
    }

    @Test
    void loginUser_nullUserLoginDto_throwsIllegalArgumentException() {
        // Expect exception
        assertThrows(IllegalArgumentException.class, () -> authenticationService.loginUser(null));

        // No interactions should happen
        verifyNoInteractions(userRepository, passwordEncoder, jwtService);
    }

    @Test
    void loginUser_blankEmail_throwsIllegalArgumentException() {
        UserLoginDto userLoginDto = new UserLoginDto("", VALID_PASSWORD);

        assertThrows(IllegalArgumentException.class, () -> authenticationService.loginUser(userLoginDto));

        verifyNoInteractions(userRepository, passwordEncoder, jwtService);
    }

    // Tests for registerUser method

    @Test
    void registerUser_successfulRegistration_returnsUserDtoWithJwtToken() {
        UserCreationDTO userCreationDTO = new UserCreationDTO(VALID_EMAIL, VALID_PASSWORD, VALID_PASSWORD);
        UserDto createdUser = UserDto.builder().userId(1L).email(VALID_EMAIL).build();

        // Mock the behavior of userService and jwtService
        when(userService.createUser(userCreationDTO)).thenReturn(createdUser);
        when(jwtService.generateJwtToken(VALID_EMAIL)).thenReturn(JWT_TOKEN);

        // Call the method to test
        UserDto result = authenticationService.registerUser(userCreationDTO);

        // Assertions
        assertNotNull(result);
        assertEquals(1L, result.getUserId());
        assertEquals(VALID_EMAIL, result.getEmail());
        assertEquals(JWT_TOKEN, result.getJwtToken());

        // Verify interactions
        verify(userService).createUser(userCreationDTO);
        verify(jwtService).generateJwtToken(VALID_EMAIL);
    }

    @Test
    void registerUser_nullUserCreationDTO_throwsIllegalArgumentException() {
        // Expect exception
        assertThrows(IllegalArgumentException.class, () -> authenticationService.registerUser(null));

        // Verify no interactions
        verifyNoInteractions(userService, jwtService);
    }

    @Test
    void registerUser_blankEmailInUserCreationDTO_throwsIllegalArgumentException() {
        UserCreationDTO userCreationDTO = new UserCreationDTO("", VALID_PASSWORD, VALID_PASSWORD);

        // Expect exception
        assertThrows(IllegalArgumentException.class, () -> authenticationService.registerUser(userCreationDTO));

        // Verify no interactions
        verifyNoInteractions(userService, jwtService);
    }

    @Test
    void registerUser_passwordMismatch_throwsInvalidPasswordException() {
        UserCreationDTO userCreationDTO = new UserCreationDTO(VALID_EMAIL, VALID_PASSWORD, INVALID_PASSWORD);

        // Expect exception
        assertThrows(InvalidPasswordException.class, () -> authenticationService.registerUser(userCreationDTO));

        // Verify no interactions
        verifyNoInteractions(userService, jwtService);
    }
}