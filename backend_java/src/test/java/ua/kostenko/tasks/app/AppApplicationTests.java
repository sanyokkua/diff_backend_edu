package ua.kostenko.tasks.app;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.MethodOrderer;
import org.junit.jupiter.api.Order;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestMethodOrder;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.annotation.Import;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.test.annotation.Rollback;
import org.springframework.test.web.servlet.MockMvc;
import ua.kostenko.tasks.app.dto.ResponseDto;
import ua.kostenko.tasks.app.dto.task.TaskCreationDTO;
import ua.kostenko.tasks.app.dto.task.TaskDto;
import ua.kostenko.tasks.app.dto.task.TaskUpdateDTO;
import ua.kostenko.tasks.app.dto.user.*;

import java.util.ArrayList;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultHandlers.print;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@AutoConfigureMockMvc
@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
@Rollback(false)
@Import(TestcontainersConfiguration.class)
@SpringBootTest
class AppApplicationTests {

    private static final ObjectMapper objectMapper = new ObjectMapper();
    private static final String BASE_API_URL = "/api/v1";
    private static final String BASE_AUTH_URL = BASE_API_URL + "/auth";
    private static final String BASE_USER_URL = BASE_API_URL + "/users";
    private static final String BASE_TASK_URL = BASE_USER_URL + "/{userId}/tasks";
    private static final String TEST_USER_1_EMAIL = "testUser1@email.com";
    private static final String TEST_USER_1_PASSWORD = "testUser1Password";
    private static final String TEST_USER_2_EMAIL = "testUser2@email.com";
    private static final String TEST_USER_2_PASSWORD = "testUser2Password";

    private static String user1Token;
    private static String user2Token;
    private static Long user1Id;
    private static Long user2Id;
    private static Long user1Task1Id;
    private static Long user1Task2Id;
    private static Long user1Task3Id;

    @Autowired
    private MockMvc mockMvc;

    private String extractToken(String responseString) throws JsonProcessingException {
        ResponseDto<UserDto> response = objectMapper.readValue(responseString, new TypeReference<>() {});
        return response.getData().getJwtToken();
    }

    private Long extractUserId(String responseString) throws JsonProcessingException {
        ResponseDto<UserDto> response = objectMapper.readValue(responseString, new TypeReference<>() {});
        return response.getData().getUserId();
    }

    private Long extractTaskId(String responseString) throws JsonProcessingException {
        ResponseDto<TaskDto> response = objectMapper.readValue(responseString, new TypeReference<>() {});
        return response.getData().getTaskId();
    }

    @Order(1)
    @Test
    void contextLoads() {
        // Verifies the application context loads without issues
    }

    @Order(2)
    @Test
    void registerUser1_POST_createsUserOne() throws Exception {
        String requestJson = objectMapper.writeValueAsString(new UserCreationDTO(TEST_USER_1_EMAIL,
                                                                                 TEST_USER_1_PASSWORD,
                                                                                 TEST_USER_1_PASSWORD));

        mockMvc.perform(post(BASE_AUTH_URL + "/register").contentType(MediaType.APPLICATION_JSON).content(requestJson))
               .andExpect(status().isCreated())
               .andExpect(jsonPath("$.statusCode").value(HttpStatus.CREATED.value()))
               .andExpect(jsonPath("$.data.userId").exists())
               .andExpect(jsonPath("$.data.email").value(TEST_USER_1_EMAIL))
               .andExpect(jsonPath("$.error").doesNotExist())
               .andDo(print());
    }

    @Order(3)
    @Test
    void registerUser2_POST_createsUserTwo() throws Exception {
        String requestJson = objectMapper.writeValueAsString(new UserCreationDTO(TEST_USER_2_EMAIL,
                                                                                 TEST_USER_2_PASSWORD,
                                                                                 TEST_USER_2_PASSWORD));

        mockMvc.perform(post(BASE_AUTH_URL + "/register").contentType(MediaType.APPLICATION_JSON).content(requestJson))
               .andExpect(status().isCreated())
               .andExpect(jsonPath("$.statusCode").value(HttpStatus.CREATED.value()))
               .andExpect(jsonPath("$.data.userId").exists())
               .andExpect(jsonPath("$.data.email").value(TEST_USER_2_EMAIL))
               .andExpect(jsonPath("$.error").doesNotExist())
               .andDo(print());
    }

    @Order(4)
    @Test
    void loginUser1_POST_logsInUserOne() throws Exception {
        var respStr = loginUser(TEST_USER_1_EMAIL, TEST_USER_1_PASSWORD);
        assertNotNull(respStr);

        user1Token = extractToken(respStr);
        user1Id = extractUserId(respStr);
    }

    @Order(5)
    @Test
    void loginUser2_POST_logsInUserTwo() throws Exception {
        var respStr = loginUser(TEST_USER_2_EMAIL, TEST_USER_2_PASSWORD);
        assertNotNull(respStr);

        user2Token = extractToken(respStr);
        user2Id = extractUserId(respStr);
    }

    private String loginUser(String email, String password) throws Exception {
        String requestJson = objectMapper.writeValueAsString(new UserLoginDto(email, password));

        return mockMvc.perform(post(BASE_AUTH_URL + "/login").contentType(MediaType.APPLICATION_JSON)
                                                             .content(requestJson))
                      .andExpect(status().isOk())
                      .andExpect(jsonPath("$.data.userId").exists())
                      .andExpect(jsonPath("$.data.jwtToken").exists())
                      .andDo(print())
                      .andReturn()
                      .getResponse()
                      .getContentAsString();
    }

    @Order(6)
    @Test
    void user1ChangePassword_PUT_changesUser1Password() throws Exception {
        String newPassword = "newPassword1";
        String requestJson =
                objectMapper.writeValueAsString(new UserUpdateDTO(TEST_USER_1_PASSWORD, newPassword, newPassword));

        mockMvc.perform(put(BASE_USER_URL + "/1/password").contentType(MediaType.APPLICATION_JSON)
                                                          .content(requestJson)
                                                          .header("Authorization", "Bearer " + user1Token))
               .andExpect(status().isOk())
               .andExpect(jsonPath("$.data.email").value(TEST_USER_1_EMAIL))
               .andDo(print());
    }

    @Order(7)
    @Test
    void loginUser1_POST_loginWithNewPassword() throws Exception {
        String newPassword = "newPassword1";
        var respStr = loginUser(TEST_USER_1_EMAIL, newPassword);
        assertNotNull(respStr);

        user1Id = extractUserId(respStr);
        user1Token = extractToken(respStr);
    }

    @Order(8)
    @Test
    void user2CreateTask_POST_createsTaskForUserTwo() throws Exception {
        createTaskForUser(user2Token, user2Id, "Task 1", "Task 1 Description");
    }

    @Order(9)
    @Test
    void user2DeleteAccount_POST_deletesUserTwo() throws Exception {
        var dto = UserDeletionDTO.builder().email(TEST_USER_2_EMAIL).currentPassword(TEST_USER_2_PASSWORD).build();
        mockMvc.perform(post(BASE_USER_URL + "/" + user2Id + "/delete").header("Authorization", "Bearer " + user2Token)
                                                                       .contentType(MediaType.APPLICATION_JSON)
                                                                       .content(objectMapper.writeValueAsString(dto)))
               .andExpect(status().isNoContent())
               .andDo(print());
    }

    @Order(10)
    @Test
    void user1CreateTask_POST_createsTaskForUserOne() throws Exception {
        user1Task1Id = createTaskForUser(user1Token, user1Id, "Task 1", "Task 1 Description");
    }

    @Order(11)
    @Test
    void user1UpdateTask_PUT_updatesTaskForUserOne() throws Exception {
        String requestJson =
                objectMapper.writeValueAsString(new TaskUpdateDTO("Updated Task Name", "Updated Task Description"));

        mockMvc.perform(put(BASE_TASK_URL + "/" + user1Task1Id, user1Id).contentType(MediaType.APPLICATION_JSON)
                                                                        .content(requestJson)
                                                                        .header("Authorization",
                                                                                "Bearer " + user1Token))
               .andExpect(status().isOk())
               .andExpect(jsonPath("$.data.name").value("Updated Task Name"))
               .andExpect(jsonPath("$.data.description").value("Updated Task Description"))
               .andDo(print());
    }

    @Order(12)
    @Test
    void user1DeleteTask_DELETE_deletesTaskForUserOne() throws Exception {
        deleteTaskForUser(user1Token, user1Id, user1Task1Id);
    }

    private Long createTaskForUser(String token, Long userId, String taskName,
                                   String taskDescription) throws Exception {
        String requestJson = objectMapper.writeValueAsString(new TaskCreationDTO(taskName, taskDescription));

        String responseString =
                mockMvc.perform(post(BASE_TASK_URL + "/", userId).contentType(MediaType.APPLICATION_JSON)
                                                                 .content(requestJson)
                                                                 .header("Authorization", "Bearer " + token))
                       .andExpect(status().isCreated())
                       .andExpect(jsonPath("$.data.taskId").exists())
                       .andExpect(jsonPath("$.data.name").value(taskName))
                       .andDo(print())
                       .andReturn()
                       .getResponse()
                       .getContentAsString();

        return extractTaskId(responseString);
    }

    private void deleteTaskForUser(String token, Long userId, Long taskId) throws Exception {
        mockMvc.perform(delete(BASE_TASK_URL + "/" + taskId, userId).header("Authorization", "Bearer " + token))
               .andExpect(status().isNoContent())
               .andDo(print());
    }

    @Order(13)
    @Test
    void user1CreateMultipleTasks_POST_createsMultipleTasksForUserOne() throws Exception {
        List<TaskCreationDTO> tasksToCreate = List.of(new TaskCreationDTO("Task 1", "Task 1 Description"),
                                                      new TaskCreationDTO("Task 2", "Task 2 Description"),
                                                      new TaskCreationDTO("Task 3", "Task 3 Description"));

        List<Long> taskIds = new ArrayList<>();
        for (TaskCreationDTO task : tasksToCreate) {
            Long taskId = createTaskForUser(user1Token, user1Id, task.getName(), task.getDescription());
            taskIds.add(taskId);
        }

        user1Task1Id = taskIds.get(0);
        user1Task2Id = taskIds.get(1);
        user1Task3Id = taskIds.get(2);
    }

    @Order(14)
    @Test
    void user1GetTasks_GET_retrievesAllTasksForUserOne() throws Exception {
        ResponseDto<List<TaskDto>> tasks = getTasksForUser(user1Token, user1Id);

        assertNotNull(tasks);
        assertEquals(3, tasks.getData().size());

        List<Long> ids = tasks.getData().stream().map(TaskDto::getTaskId).toList();
        assertTrue(ids.contains(user1Task1Id));
        assertTrue(ids.contains(user1Task2Id));
        assertTrue(ids.contains(user1Task3Id));
    }

    @Order(15)
    @Test
    void user1DeleteTaskAgain_DELETE_deletesTaskForUserOne() throws Exception {
        deleteTaskForUser(user1Token, user1Id, user1Task1Id);
    }

    @Order(16)
    @Test
    void user1GetTasksAgain_GET_retrievesAllTasksForUserOne() throws Exception {
        ResponseDto<List<TaskDto>> tasks = getTasksForUser(user1Token, user1Id);

        assertNotNull(tasks);
        assertEquals(2, tasks.getData().size());

        List<Long> ids = tasks.getData().stream().map(TaskDto::getTaskId).toList();
        assertFalse(ids.contains(user1Task1Id));
        assertTrue(ids.contains(user1Task2Id));
        assertTrue(ids.contains(user1Task3Id));
    }

    @Order(17)
    @Test
    void user1DeleteTaskNotExisting_DELETE_deletesTaskForUserOne() throws Exception {
        mockMvc.perform(delete(BASE_TASK_URL + "/" + 99999L, user1Id).header("Authorization", "Bearer " + user1Token))
               .andExpect(status().isNotFound())
               .andExpect(jsonPath("$.statusCode").value(HttpStatus.NOT_FOUND.value()))
               .andExpect(jsonPath("$.statusMessage").value(HttpStatus.NOT_FOUND.name()))
               .andExpect(jsonPath("$.data").doesNotExist())
               .andExpect(jsonPath("$.error").exists())
               .andDo(print());
    }

    // Helper method to retrieve all tasks for a user
    private ResponseDto<List<TaskDto>> getTasksForUser(String token, Long userId) throws Exception {
        String responseString =
                mockMvc.perform(get(BASE_TASK_URL + "/", userId).header("Authorization", "Bearer " + token))
                       .andExpect(status().isOk())
                       .andExpect(jsonPath("$.statusCode").value(HttpStatus.OK.value()))
                       .andExpect(jsonPath("$.statusMessage").value(HttpStatus.OK.name()))
                       .andExpect(jsonPath("$.data").exists())
                       .andExpect(jsonPath("$.data").isArray())
                       .andDo(print())
                       .andReturn()
                       .getResponse()
                       .getContentAsString();

        return objectMapper.readValue(responseString, new TypeReference<>() {});
    }

    @Order(18)
    @Test
    void user1CreateEmptyTask_POST_failsToCreateEmptyTask() throws Exception {
        var emptyTask = new TaskCreationDTO("", "");
        var emptyTaskReq = objectMapper.writeValueAsString(emptyTask);

        mockMvc.perform(post(BASE_TASK_URL + "/", user1Id).contentType("application/json")
                                                          .content(emptyTaskReq)
                                                          .header("Authorization", "Bearer " + user1Token))
               .andExpect(status().isBadRequest())
               .andExpect(jsonPath("$.statusCode").value(HttpStatus.BAD_REQUEST.value()))
               .andExpect(jsonPath("$.statusMessage").value(HttpStatus.BAD_REQUEST.name()))
               .andExpect(jsonPath("$.data").doesNotExist())
               .andExpect(jsonPath("$.error").exists())
               .andDo(print());
    }

    @Order(19)
    @Test
    void unauthorizedUserCreateTask_POST_failsForUnauthorizedUser() throws Exception {
        var task = new TaskCreationDTO("Unauthorized Task", "Should fail");
        var taskReq = objectMapper.writeValueAsString(task);

        mockMvc.perform(post(BASE_TASK_URL + "/", user1Id).contentType("application/json").content(taskReq))
               .andExpect(status().isUnauthorized())
               .andExpect(jsonPath("$.statusCode").value(HttpStatus.UNAUTHORIZED.value()))
               .andExpect(jsonPath("$.statusMessage").value(HttpStatus.UNAUTHORIZED.name()))
               .andExpect(jsonPath("$.data").doesNotExist())
               .andExpect(jsonPath("$.error").exists())
               .andDo(print());
    }

    @Order(20)
    @Test
    void user2DeleteTaskOfUser1_DELETE_failsForUnauthorizedDeletion() throws Exception {
        mockMvc.perform(delete(BASE_TASK_URL + "/" + user1Task2Id, user2Id).header("Authorization",
                                                                                   "Bearer " + user2Token))
               .andExpect(status().isUnauthorized())
               .andExpect(jsonPath("$.statusCode").value(HttpStatus.UNAUTHORIZED.value()))
               .andExpect(jsonPath("$.statusMessage").value(HttpStatus.UNAUTHORIZED.name()))
               .andExpect(jsonPath("$.data").doesNotExist())
               .andExpect(jsonPath("$.error").exists())
               .andDo(print());
    }

    @Order(21)
    @Test
    void user1UpdateTaskWithInvalidData_PUT_failsToUpdate() throws Exception {
        var invalidTask = new TaskCreationDTO("", "Updated Description");
        var invalidTaskReq = objectMapper.writeValueAsString(invalidTask);

        mockMvc.perform(put(BASE_TASK_URL + "/" + user1Task2Id, user1Id).contentType("application/json")
                                                                        .content(invalidTaskReq)
                                                                        .header("Authorization",
                                                                                "Bearer " + user1Token))
               .andExpect(status().isBadRequest())
               .andExpect(jsonPath("$.statusCode").value(HttpStatus.BAD_REQUEST.value()))
               .andExpect(jsonPath("$.statusMessage").value(HttpStatus.BAD_REQUEST.name()))
               .andExpect(jsonPath("$.data").doesNotExist())
               .andExpect(jsonPath("$.error").exists())
               .andDo(print());
    }

    @Order(22)
    @Test
    void user1GetSingleTaskById_GET_retrievesSingleTask() throws Exception {
        mockMvc.perform(get(BASE_TASK_URL + "/" + user1Task2Id, user1Id).header("Authorization",
                                                                                "Bearer " + user1Token))
               .andExpect(status().isOk())
               .andExpect(jsonPath("$.statusCode").value(HttpStatus.OK.value()))
               .andExpect(jsonPath("$.statusMessage").value(HttpStatus.OK.name()))
               .andExpect(jsonPath("$.data.taskId").value(user1Task2Id))
               .andExpect(jsonPath("$.data.name").exists())
               .andExpect(jsonPath("$.data.description").exists())
               .andExpect(jsonPath("$.error").doesNotExist())
               .andDo(print());
    }

}
