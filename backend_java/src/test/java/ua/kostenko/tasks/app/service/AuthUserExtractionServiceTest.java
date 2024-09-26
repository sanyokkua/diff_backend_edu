package ua.kostenko.tasks.app.service;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.security.authentication.AuthenticationCredentialsNotFoundException;
import org.springframework.security.authentication.InsufficientAuthenticationException;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContext;
import org.springframework.security.core.context.SecurityContextHolder;
import ua.kostenko.tasks.app.config.custom.UserAuthentication;
import ua.kostenko.tasks.app.dto.user.UserDto;
import ua.kostenko.tasks.app.entity.User;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class AuthUserExtractionServiceTest {

    private static final Long USER_ID = 1L;
    private static final String USER_EMAIL = "test@example.com";
    private static final String JWT_TOKEN = "jwt-token";
    @InjectMocks
    private AuthUserExtractionService authUserExtractionService;
    @Mock
    private SecurityContext securityContext;
    @Mock
    private Authentication authentication;
    @Mock
    private UserAuthentication userAuthentication;
    @Mock
    private User user;

    @BeforeEach
    void setUp() {
        MockitoAnnotations.openMocks(this);
    }

    @Test
    void getUserFromAuthContext_Success() {
        // Arrange
        when(securityContext.getAuthentication()).thenReturn(userAuthentication);
        when(userAuthentication.getPrincipal()).thenReturn(user);
        when(user.getUserId()).thenReturn(USER_ID);
        when(user.getEmail()).thenReturn(USER_EMAIL);
        when(userAuthentication.getJwtToken()).thenReturn(JWT_TOKEN);

        SecurityContextHolder.setContext(securityContext);

        // Act
        UserDto result = authUserExtractionService.getUserFromAuthContext();

        // Assert
        assertNotNull(result);
        assertEquals(USER_ID, result.getUserId());
        assertEquals(USER_EMAIL, result.getEmail());
        assertEquals(JWT_TOKEN, result.getJwtToken());
    }

    @Test
    void getUserFromAuthContext_AuthenticationIsNull_ShouldThrowException() {
        // Arrange
        when(securityContext.getAuthentication()).thenReturn(null);
        SecurityContextHolder.setContext(securityContext);

        // Act & Assert
        AuthenticationCredentialsNotFoundException exception =
                assertThrows(AuthenticationCredentialsNotFoundException.class,
                             () -> authUserExtractionService.getUserFromAuthContext());

        assertEquals("User is not authenticated", exception.getMessage());
        verify(securityContext).getAuthentication();
    }

    @Test
    void getUserFromAuthContext_AuthenticationNotUserAuthentication_ShouldThrowException() {
        // Arrange
        when(securityContext.getAuthentication()).thenReturn(authentication); // Not UserAuthentication
        SecurityContextHolder.setContext(securityContext);

        // Act & Assert
        AuthenticationCredentialsNotFoundException exception =
                assertThrows(AuthenticationCredentialsNotFoundException.class,
                             () -> authUserExtractionService.getUserFromAuthContext());

        assertEquals("User is not authenticated", exception.getMessage());
        verify(securityContext).getAuthentication();
    }

    @Test
    void extractUserDto_PrincipalIsNull_ShouldThrowException() {
        // Arrange
        when(securityContext.getAuthentication()).thenReturn(userAuthentication);
        when(userAuthentication.getPrincipal()).thenReturn(null);
        SecurityContextHolder.setContext(securityContext);
        // Act & Assert
        InsufficientAuthenticationException exception = assertThrows(InsufficientAuthenticationException.class,
                                                                     () -> authUserExtractionService.getUserFromAuthContext());

        assertEquals("User is not set to Authentication", exception.getMessage());
        verify(userAuthentication).getPrincipal();
    }

    @Test
    void extractUserDto_UserDetailsAreNull_ShouldReturnIncompleteUserDto() {
        // Arrange
        when(securityContext.getAuthentication()).thenReturn(userAuthentication);
        when(userAuthentication.getPrincipal()).thenReturn(user);
        when(user.getUserId()).thenReturn(null);  // Null user ID
        when(user.getEmail()).thenReturn(USER_EMAIL);
        when(userAuthentication.getJwtToken()).thenReturn(JWT_TOKEN);

        SecurityContextHolder.setContext(securityContext);

        // Act
        UserDto result = authUserExtractionService.getUserFromAuthContext();

        // Assert
        assertNotNull(result);
        assertNull(result.getUserId());
        assertEquals(USER_EMAIL, result.getEmail());
        assertEquals(JWT_TOKEN, result.getJwtToken());
    }

    @Test
    void extractUserDto_JwtTokenIsEmpty_ShouldStillReturnUserDto() {
        // Arrange
        when(securityContext.getAuthentication()).thenReturn(userAuthentication);
        when(userAuthentication.getPrincipal()).thenReturn(user);
        when(user.getUserId()).thenReturn(USER_ID);
        when(user.getEmail()).thenReturn(USER_EMAIL);
        when(userAuthentication.getJwtToken()).thenReturn("");  // Empty JWT token

        SecurityContextHolder.setContext(securityContext);

        // Act
        UserDto result = authUserExtractionService.getUserFromAuthContext();

        // Assert
        assertNotNull(result);
        assertEquals(USER_ID, result.getUserId());
        assertEquals(USER_EMAIL, result.getEmail());
        assertEquals("", result.getJwtToken());  // Empty JWT token
    }

}