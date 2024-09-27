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
import ua.kostenko.tasks.app.dto.user.UserDeletionDTO;
import ua.kostenko.tasks.app.dto.user.UserDto;
import ua.kostenko.tasks.app.dto.user.UserUpdateDTO;
import ua.kostenko.tasks.app.exception.InvalidPasswordException;
import ua.kostenko.tasks.app.repository.UserRepository;
import ua.kostenko.tasks.app.service.AuthUserExtractionService;
import ua.kostenko.tasks.app.service.JwtService;
import ua.kostenko.tasks.app.service.UserService;

import static org.mockito.Mockito.doThrow;
import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultHandlers.print;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@WebMvcTest(UserController.class)
@AutoConfigureMockMvc(addFilters = false)
class UserControllerTest {

    private static final String BASE_URL = "/api/v1/users";
    private static final Long VALID_USER_ID = 1L;
    private static final String VALID_EMAIL = "valid@email.com";
    private static final String VALID_PASSWORD = "testPassword";

    private final ObjectMapper objectMapper = new ObjectMapper();
    // Mock UserDto to use in tests
    private final UserDto mockUserDto = UserDto.builder().userId(VALID_USER_ID).email(VALID_EMAIL).build();
    @MockBean
    private UserService userService;
    @MockBean
    private JwtService jwtService;
    @MockBean
    private UserRepository userRepository;
    @MockBean
    private AuthUserExtractionService userExtractionService;
    @Autowired
    private MockMvc mockMvc;

    // Test Methods
    @Test
    @DisplayName("Fetch User Details - Valid Request")
    void getUserById_validRequest_shouldReturnUserDto() throws Exception {
        UserDto userDto = UserDto.builder().userId(VALID_USER_ID).email(VALID_EMAIL).build();

        when(userExtractionService.getUserFromAuthContext()).thenReturn(userDto);

        performGetRequest("/" + VALID_USER_ID).andExpect(status().isOk())
                                              .andExpect(jsonPath("$.statusCode").value(200))
                                              .andExpect(jsonPath("$.statusMessage").value("OK"))
                                              .andExpect(jsonPath("$.data.email").value(VALID_EMAIL))
                                              .andExpect(jsonPath("$.error").doesNotExist())
                                              .andDo(print());
    }

    @Test
    @DisplayName("Update User Password - Valid Request")
    void updateUserPassword_validRequest_shouldReturnUpdatedUserDto() throws Exception {
        UserUpdateDTO updateDTO = UserUpdateDTO.builder()
                                               .currentPassword(VALID_PASSWORD)
                                               .newPassword("newPassword")
                                               .newPasswordConfirmation("newPassword")
                                               .build();

        UserDto updatedUserDto = UserDto.builder().userId(VALID_USER_ID).email(VALID_EMAIL).build();

        when(userExtractionService.getUserFromAuthContext()).thenReturn(updatedUserDto);
        when(userService.updateUserPassword(VALID_USER_ID, updateDTO)).thenReturn(updatedUserDto);

        performPutRequest("/" + VALID_USER_ID + "/password", updateDTO).andExpect(status().isOk())
                                                                       .andExpect(jsonPath("$.statusCode").value(200))
                                                                       .andExpect(jsonPath("$.statusMessage").value("OK"))
                                                                       .andExpect(jsonPath("$.data.email").value(
                                                                               VALID_EMAIL))
                                                                       .andExpect(jsonPath("$.error").doesNotExist())
                                                                       .andDo(print());
    }

    @Test
    @DisplayName("Update User Password - Invalid Request")
    void updateUserPassword_invalidRequest_shouldReturnException() throws Exception {
        UserUpdateDTO updateDTO = UserUpdateDTO.builder()
                                               .currentPassword(VALID_PASSWORD)
                                               .newPassword("newPassword")
                                               .newPasswordConfirmation("differentPassword")
                                               .build();

        when(userExtractionService.getUserFromAuthContext()).thenReturn(mockUserDto);
        when(userService.updateUserPassword(VALID_USER_ID, updateDTO)).thenThrow(new InvalidPasswordException(
                "Passwords do not match"));

        performPutRequest("/" + VALID_USER_ID + "/password", updateDTO).andExpect(status().isBadRequest())
                                                                       .andExpect(jsonPath("$.statusCode").value(400))
                                                                       .andExpect(jsonPath("$.statusMessage").value(
                                                                               "BAD_REQUEST"))
                                                                       .andExpect(jsonPath("$.error").value(
                                                                               "InvalidPasswordException: Passwords do not match"))
                                                                       .andDo(print());
    }

    @Test
    @DisplayName("Delete User - Valid Request")
    void deleteUser_validRequest_shouldReturnNoContent() throws Exception {
        UserDeletionDTO deletionDTO = new UserDeletionDTO();
        deletionDTO.setCurrentPassword(VALID_PASSWORD);

        when(userExtractionService.getUserFromAuthContext()).thenReturn(mockUserDto);

        performPostRequest("/" + VALID_USER_ID + "/delete", deletionDTO).andExpect(status().isNoContent())
                                                                        .andDo(print());
    }

    @Test
    @DisplayName("Delete User - Invalid Request")
    void deleteUser_invalidRequest_shouldReturnException() throws Exception {
        UserDeletionDTO deletionDTO = new UserDeletionDTO();
        deletionDTO.setCurrentPassword(VALID_PASSWORD);

        when(userExtractionService.getUserFromAuthContext()).thenReturn(mockUserDto);
        doThrow(new InvalidPasswordException("Invalid credentials")).when(userService)
                                                                    .deleteUser(VALID_USER_ID, deletionDTO);

        performPostRequest("/" + VALID_USER_ID + "/delete", deletionDTO).andExpect(status().isBadRequest())
                                                                        .andExpect(jsonPath("$.statusCode").value(400))
                                                                        .andExpect(jsonPath("$.statusMessage").value(
                                                                                "BAD_REQUEST"))
                                                                        .andExpect(jsonPath("$.error").value(
                                                                                "InvalidPasswordException: Invalid credentials"))
                                                                        .andDo(print());
    }

    // Helper Methods
    private ResultActions performGetRequest(String endpoint) throws Exception {
        return mockMvc.perform(get(BASE_URL + endpoint).contentType(MediaType.APPLICATION_JSON));
    }

    private ResultActions performPutRequest(String endpoint, Object requestDto) throws Exception {
        return mockMvc.perform(put(BASE_URL + endpoint).contentType(MediaType.APPLICATION_JSON)
                                                       .content(objectMapper.writeValueAsString(requestDto)));
    }

    private ResultActions performPostRequest(String endpoint, Object requestDto) throws Exception {
        return mockMvc.perform(post(BASE_URL + endpoint).contentType(MediaType.APPLICATION_JSON)
                                                        .content(objectMapper.writeValueAsString(requestDto)));
    }
}