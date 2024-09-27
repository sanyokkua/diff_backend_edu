package ua.kostenko.tasks.app.config.custom;

import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.servlet.ServletOutputStream;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.http.HttpStatus;
import org.springframework.security.access.AccessDeniedException;
import ua.kostenko.tasks.app.utility.ResponseDtoUtils;

import java.io.IOException;

import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class CustomAccessDeniedHandlerTest {

    private CustomAccessDeniedHandler customAccessDeniedHandler;
    private ObjectMapper objectMapper;
    private HttpServletRequest request;
    private HttpServletResponse response;
    private AccessDeniedException accessDeniedException;

    @BeforeEach
    void setUp() throws IOException {
        objectMapper = mock(ObjectMapper.class);
        request = mock(HttpServletRequest.class);
        response = mock(HttpServletResponse.class);
        ServletOutputStream servletOutputStream = mock(ServletOutputStream.class);
        when(response.getOutputStream()).thenReturn(servletOutputStream);

        accessDeniedException = new AccessDeniedException("You do not have permission to access this resource");
        customAccessDeniedHandler = new CustomAccessDeniedHandler(objectMapper);
    }

    @Test
    void testHandle_Success() throws IOException {
        // Arrange
        String requestURI = "/test/url";
        when(request.getRequestURI()).thenReturn(requestURI);
        var errorResponse = ResponseDtoUtils.createErrorResponseBody(null, HttpStatus.FORBIDDEN, accessDeniedException);
        String jsonResponse = "{\"error\":\"forbidden\"}";

        when(objectMapper.writeValueAsString(errorResponse)).thenReturn(jsonResponse);

        // Act
        customAccessDeniedHandler.handle(request, response, accessDeniedException);

        // Assert
        verify(response).setContentType("application/json");
        verify(response).setStatus(HttpServletResponse.SC_FORBIDDEN);
        verify(response.getOutputStream()).println(jsonResponse);
    }

}