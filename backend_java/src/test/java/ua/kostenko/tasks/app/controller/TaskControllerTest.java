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
import ua.kostenko.tasks.app.dto.task.TaskCreationDTO;
import ua.kostenko.tasks.app.dto.task.TaskDto;
import ua.kostenko.tasks.app.dto.task.TaskUpdateDTO;
import ua.kostenko.tasks.app.dto.user.UserDto;
import ua.kostenko.tasks.app.repository.UserRepository;
import ua.kostenko.tasks.app.service.AuthUserExtractionService;
import ua.kostenko.tasks.app.service.JwtService;
import ua.kostenko.tasks.app.service.TaskService;

import java.util.List;
import java.util.Optional;

import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultHandlers.print;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@WebMvcTest(TaskController.class)
@AutoConfigureMockMvc(addFilters = false)
class TaskControllerTest {

    private static final String BASE_URL = "/api/v1/users/{userId}/tasks";
    private static final Long VALID_USER_ID = 1L;
    private static final Long VALID_TASK_ID = 1L;
    private static final String VALID_TASK_NAME = "Test Task";

    private final ObjectMapper objectMapper = new ObjectMapper();
    // Mock UserDto for tests
    private final UserDto mockUserDto = UserDto.builder().userId(VALID_USER_ID).email("valid@email.com").build();
    @MockBean
    private TaskService taskService;
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
    @DisplayName("Create Task - Valid Request")
    void createTask_validRequest_shouldReturnCreatedTaskDto() throws Exception {
        TaskCreationDTO taskCreationDTO = new TaskCreationDTO();
        taskCreationDTO.setName(VALID_TASK_NAME);

        TaskDto createdTaskDto = TaskDto.builder().taskId(VALID_TASK_ID).name(VALID_TASK_NAME).build();

        when(userExtractionService.getUserFromAuthContext()).thenReturn(mockUserDto);
        when(taskService.createTask(VALID_USER_ID, taskCreationDTO)).thenReturn(createdTaskDto);

        performPostRequest(taskCreationDTO).andExpect(status().isCreated())
                                           .andExpect(jsonPath("$.statusCode").value(201))
                                           .andExpect(jsonPath("$.statusMessage").value("CREATED"))
                                           .andExpect(jsonPath("$.data.name").value(VALID_TASK_NAME))
                                           .andExpect(jsonPath("$.error").doesNotExist())
                                           .andDo(print());
    }

    @Test
    @DisplayName("Get Task by ID - Valid Request")
    void getTaskById_validRequest_shouldReturnTaskDto() throws Exception {
        TaskDto taskDto = TaskDto.builder().taskId(VALID_TASK_ID).name(VALID_TASK_NAME).build();

        when(userExtractionService.getUserFromAuthContext()).thenReturn(mockUserDto);
        when(taskService.getTaskByUserIdAndTaskId(VALID_USER_ID, VALID_TASK_ID)).thenReturn(Optional.of(taskDto));

        performGetRequest("/" + VALID_TASK_ID).andExpect(status().isOk())
                                              .andExpect(jsonPath("$.statusCode").value(200))
                                              .andExpect(jsonPath("$.statusMessage").value("OK"))
                                              .andExpect(jsonPath("$.data.name").value(VALID_TASK_NAME))
                                              .andExpect(jsonPath("$.error").doesNotExist())
                                              .andDo(print());
    }

    @Test
    @DisplayName("Get Task by ID - Not Found Request")
    void getTaskById_notFound_shouldReturnError() throws Exception {
        when(userExtractionService.getUserFromAuthContext()).thenReturn(mockUserDto);
        when(taskService.getTaskByUserIdAndTaskId(VALID_USER_ID, VALID_TASK_ID)).thenReturn(Optional.empty());

        performGetRequest("/" + VALID_TASK_ID).andExpect(status().isNotFound())
                                              .andExpect(jsonPath("$.statusCode").value(404))
                                              .andExpect(jsonPath("$.statusMessage").value("NOT_FOUND"))
                                              .andExpect(jsonPath("$.error").doesNotExist())
                                              .andDo(print());
    }

    @Test
    @DisplayName("Update Task - Valid Request")
    void updateTask_validRequest_shouldReturnUpdatedTaskDto() throws Exception {
        TaskUpdateDTO taskUpdateDTO = new TaskUpdateDTO();
        taskUpdateDTO.setName("Updated Task Name");

        TaskDto updatedTaskDto = TaskDto.builder().taskId(VALID_TASK_ID).name("Updated Task Name").build();

        when(userExtractionService.getUserFromAuthContext()).thenReturn(mockUserDto);
        when(taskService.updateTask(VALID_USER_ID, VALID_TASK_ID, taskUpdateDTO)).thenReturn(updatedTaskDto);

        performPutRequest("/" + VALID_TASK_ID, taskUpdateDTO).andExpect(status().isOk())
                                                             .andExpect(jsonPath("$.statusCode").value(200))
                                                             .andExpect(jsonPath("$.statusMessage").value("OK"))
                                                             .andExpect(jsonPath("$.data.name").value(
                                                                     "Updated Task Name"))
                                                             .andExpect(jsonPath("$.error").doesNotExist())
                                                             .andDo(print());
    }

    @Test
    @DisplayName("Delete Task - Valid Request")
    void deleteTask_validRequest_shouldReturnNoContent() throws Exception {
        when(userExtractionService.getUserFromAuthContext()).thenReturn(mockUserDto);

        performDeleteRequest("/" + VALID_TASK_ID).andExpect(status().isNoContent()).andDo(print());

        verify(taskService).deleteTask(VALID_USER_ID, VALID_TASK_ID);
    }

    @Test
    @DisplayName("Get All Tasks - Valid Request")
    void getAllTasks_validRequest_shouldReturnTaskList() throws Exception {
        TaskDto taskDto1 = TaskDto.builder().taskId(1L).name("Task 1").build();
        TaskDto taskDto2 = TaskDto.builder().taskId(2L).name("Task 2").build();
        List<TaskDto> tasks = List.of(taskDto1, taskDto2);

        when(userExtractionService.getUserFromAuthContext()).thenReturn(mockUserDto);
        when(taskService.getAllTasksForUser(VALID_USER_ID)).thenReturn(tasks);

        performGetRequest("/").andExpect(status().isOk())
                              .andExpect(jsonPath("$.statusCode").value(200))
                              .andExpect(jsonPath("$.statusMessage").value("OK"))
                              .andExpect(jsonPath("$.data").isArray())
                              .andExpect(jsonPath("$.data.length()").value(2))
                              .andExpect(jsonPath("$.error").doesNotExist())
                              .andDo(print());
    }

    // Helper Methods
    private ResultActions performPostRequest(Object requestDto) throws Exception {
        return mockMvc.perform(post(BASE_URL + "/",
                                    TaskControllerTest.VALID_USER_ID).contentType(MediaType.APPLICATION_JSON)
                                                                     .content(objectMapper.writeValueAsString(requestDto)));
    }

    private ResultActions performGetRequest(String endpoint) throws Exception {
        return mockMvc.perform(get(BASE_URL + endpoint,
                                   TaskControllerTest.VALID_USER_ID).contentType(MediaType.APPLICATION_JSON));
    }

    private ResultActions performPutRequest(String endpoint, Object requestDto) throws Exception {
        return mockMvc.perform(put(BASE_URL + endpoint,
                                   TaskControllerTest.VALID_USER_ID).contentType(MediaType.APPLICATION_JSON)
                                                                    .content(objectMapper.writeValueAsString(requestDto)));
    }

    private ResultActions performDeleteRequest(String endpoint) throws Exception {
        return mockMvc.perform(delete(BASE_URL + endpoint,
                                      TaskControllerTest.VALID_USER_ID).contentType(MediaType.APPLICATION_JSON));
    }
}