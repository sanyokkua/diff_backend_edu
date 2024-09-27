package ua.kostenko.tasks.app.config.custom;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.security.access.AccessDeniedException;
import org.springframework.security.authentication.AuthenticationCredentialsNotFoundException;
import org.springframework.security.authentication.InsufficientAuthenticationException;
import org.springframework.test.web.servlet.MockMvc;
import ua.kostenko.tasks.app.controller.AuthController;
import ua.kostenko.tasks.app.dto.user.UserLoginDto;
import ua.kostenko.tasks.app.exception.*;
import ua.kostenko.tasks.app.repository.UserRepository;
import ua.kostenko.tasks.app.service.AuthenticationService;
import ua.kostenko.tasks.app.service.JwtService;

import java.util.stream.Stream;

import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultHandlers.print;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@WebMvcTest(AuthController.class)// it doesn't matter what controller is used
@AutoConfigureMockMvc(addFilters = false)
class CustomRestExceptionHandlerTest {

    private static final String BASE_URL = "/api/v1/auth/login";
    private static final String VALID_EMAIL = "valid@email.com";
    private static final String VALID_PASSWORD = "testPassword";
    private static final String ERROR = "TestError";
    private final ObjectMapper objectMapper = new ObjectMapper();
    @MockBean
    private JwtService jwtService;
    @MockBean
    private UserRepository userRepository;
    @MockBean
    private AuthenticationService authenticationService;
    @Autowired
    private MockMvc mockMvc;

    public static Stream<Arguments> parameters() {
        return Stream.of(Arguments.of(new IllegalArgumentException(ERROR),
                                      HttpStatus.BAD_REQUEST,
                                      String.format("%s: %s", "IllegalArgumentException", ERROR)),
                         Arguments.of(new InvalidEmailFormatException(ERROR),
                                      HttpStatus.BAD_REQUEST,
                                      String.format("%s: %s", "InvalidEmailFormatException", ERROR)),
                         Arguments.of(new InvalidPasswordException(ERROR),
                                      HttpStatus.BAD_REQUEST,
                                      String.format("%s: %s", "InvalidPasswordException", ERROR)),
                         Arguments.of(new InvalidJwtTokenException(ERROR),
                                      HttpStatus.BAD_REQUEST,
                                      String.format("%s: %s", "InvalidJwtTokenException", ERROR)),
                         Arguments.of(new AuthenticationCredentialsNotFoundException(ERROR),
                                      HttpStatus.UNAUTHORIZED,
                                      String.format("%s: %s", "AuthenticationCredentialsNotFoundException", ERROR)),
                         Arguments.of(new InsufficientAuthenticationException(ERROR),
                                      HttpStatus.UNAUTHORIZED,
                                      String.format("%s: %s", "InsufficientAuthenticationException", ERROR)),
                         Arguments.of(new AccessDeniedException(ERROR),
                                      HttpStatus.FORBIDDEN,
                                      String.format("%s: %s", "AccessDeniedException", ERROR)),
                         Arguments.of(new TaskNotFoundException(ERROR),
                                      HttpStatus.NOT_FOUND,
                                      String.format("%s: %s", "TaskNotFoundException", ERROR)),
                         Arguments.of(new EmailAlreadyExistsException(ERROR),
                                      HttpStatus.CONFLICT,
                                      String.format("%s: %s", "EmailAlreadyExistsException", ERROR)),
                         Arguments.of(new TaskAlreadyExistsException(ERROR),
                                      HttpStatus.CONFLICT,
                                      String.format("%s: %s", "TaskAlreadyExistsException", ERROR)),
                         Arguments.of(new RuntimeException(ERROR),
                                      HttpStatus.INTERNAL_SERVER_ERROR,
                                      String.format("%s: %s", "RuntimeException", ERROR)));
    }

    @ParameterizedTest
    @MethodSource("parameters")
    void loginUser_invalidRequest_shouldReturnError(Exception ex, HttpStatus status, String msg) throws Exception {
        UserLoginDto requestDto = UserLoginDto.builder().email(VALID_EMAIL).password(VALID_PASSWORD).build();
        when(authenticationService.loginUser(requestDto)).thenThrow(ex);

        mockMvc.perform(post(BASE_URL).contentType(MediaType.APPLICATION_JSON)
                                      .content(objectMapper.writeValueAsString(requestDto)))
               .andExpect(status().is(status.value()))
               .andExpect(jsonPath("$.statusCode").value(status.value()))
               .andExpect(jsonPath("$.statusMessage").value(status.name()))
               .andExpect(jsonPath("$.data").doesNotExist())
               .andExpect(jsonPath("$.error").exists())
               .andExpect(jsonPath("$.error").value(msg))
               .andDo(print());
    }

}