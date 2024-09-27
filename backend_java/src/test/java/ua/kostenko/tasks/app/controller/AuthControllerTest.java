package ua.kostenko.tasks.app.controller;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.ResultActions;
import ua.kostenko.tasks.app.dto.user.UserCreationDTO;
import ua.kostenko.tasks.app.dto.user.UserDto;
import ua.kostenko.tasks.app.dto.user.UserLoginDto;
import ua.kostenko.tasks.app.exception.InvalidPasswordException;
import ua.kostenko.tasks.app.repository.UserRepository;
import ua.kostenko.tasks.app.service.AuthenticationService;
import ua.kostenko.tasks.app.service.JwtService;

import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultHandlers.print;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@WebMvcTest(AuthController.class)
@AutoConfigureMockMvc(addFilters = false)
class AuthControllerTest {

    private static final String BASE_URL = "/api/v1/auth";
    private static final String VALID_EMAIL = "valid@email.com";
    private static final String VALID_PASSWORD = "testPassword";
    private final ObjectMapper objectMapper = new ObjectMapper();

    @MockBean
    private JwtService jwtService;
    @MockBean
    private UserRepository userRepository;
    @MockBean
    private AuthenticationService authenticationService;
    @Autowired
    private MockMvc mockMvc;

    // Test Methods
    @Test
    @DisplayName("Login User - Valid Request")
    void loginUser_validRequest_shouldReturnUserDto() throws Exception {
        UserLoginDto requestDto = createLoginRequest();
        UserDto userDto = UserDto.builder().userId(1L).email(requestDto.getEmail()).jwtToken("testJwtToken").build();

        when(authenticationService.loginUser(requestDto)).thenReturn(userDto);

        performPostRequest("/login", requestDto).andExpect(status().isOk())
                                                .andExpect(jsonPath("$.statusCode").value(200))
                                                .andExpect(jsonPath("$.statusMessage").value("OK"))
                                                .andExpect(jsonPath("$.data.email").value(requestDto.getEmail()))
                                                .andExpect(jsonPath("$.data.jwtToken").value("testJwtToken"))
                                                .andExpect(jsonPath("$.error").doesNotExist())
                                                .andDo(print());
    }

    @Test
    @DisplayName("Login User - Invalid Request")
    void loginUser_invalidRequest_shouldReturnException() throws Exception {
        UserLoginDto requestDto = createLoginRequest();

        when(authenticationService.loginUser(requestDto)).thenThrow(new IllegalArgumentException("User not found"));

        performPostRequest("/login", requestDto).andExpect(status().isBadRequest())
                                                .andExpect(jsonPath("$.statusCode").value(400))
                                                .andExpect(jsonPath("$.statusMessage").value("BAD_REQUEST"))
                                                .andExpect(jsonPath("$.error").value(
                                                        "IllegalArgumentException: User not found"))
                                                .andDo(print());
    }

    @Test
    @DisplayName("Register User - Valid Request")
    void registerUser_validRequest_shouldReturnUserDto() throws Exception {
        UserCreationDTO requestDto = createRegisterRequest();
        UserDto userDto = UserDto.builder().userId(1L).email(requestDto.getEmail()).jwtToken("testJwtToken").build();

        when(authenticationService.registerUser(requestDto)).thenReturn(userDto);

        performPostRequest("/register", requestDto).andExpect(status().isCreated())
                                                   .andExpect(jsonPath("$.statusCode").value(201))
                                                   .andExpect(jsonPath("$.statusMessage").value("CREATED"))
                                                   .andExpect(jsonPath("$.data.email").value(requestDto.getEmail()))
                                                   .andExpect(jsonPath("$.data.jwtToken").value("testJwtToken"))
                                                   .andExpect(jsonPath("$.error").doesNotExist())
                                                   .andDo(print());
    }

    @Test
    @DisplayName("Register User - Invalid Request")
    void registerUser_invalidRequest_shouldReturnException() throws Exception {
        UserCreationDTO requestDto = createRegisterRequest();

        when(authenticationService.registerUser(requestDto)).thenThrow(new InvalidPasswordException(
                "Invalid credentials"));

        performPostRequest("/register", requestDto).andExpect(status().isBadRequest())
                                                   .andExpect(jsonPath("$.statusCode").value(400))
                                                   .andExpect(jsonPath("$.statusMessage").value("BAD_REQUEST"))
                                                   .andExpect(jsonPath("$.error").value(
                                                           "InvalidPasswordException: Invalid credentials"))
                                                   .andDo(print());
    }

    // Helper Methods
    private ResultActions performPostRequest(String endpoint, Object requestDto) throws Exception {
        return mockMvc.perform(post(BASE_URL + endpoint).contentType(MediaType.APPLICATION_JSON)
                                                        .content(objectMapper.writeValueAsString(requestDto)));
    }

    private UserLoginDto createLoginRequest() {
        return UserLoginDto.builder().email(VALID_EMAIL).password(VALID_PASSWORD).build();
    }

    private UserCreationDTO createRegisterRequest() {
        return UserCreationDTO.builder()
                              .email(VALID_EMAIL)
                              .password(VALID_PASSWORD)
                              .passwordConfirmation(VALID_PASSWORD)
                              .build();
    }
}