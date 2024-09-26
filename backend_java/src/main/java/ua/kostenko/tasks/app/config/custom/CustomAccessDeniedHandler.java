package ua.kostenko.tasks.app.config.custom;

import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.security.access.AccessDeniedException;
import org.springframework.security.web.access.AccessDeniedHandler;
import org.springframework.stereotype.Component;
import ua.kostenko.tasks.app.utility.ResponseDtoUtils;

import java.io.IOException;

/**
 * Custom handler for managing access denial in Spring Security.
 * <p>
 * When a user attempts to access a resource they do not have permission for, this handler responds with a
 * custom JSON response, detailing the error and a 403 Forbidden status code.
 * </p>
 */
@Slf4j
@Component
@RequiredArgsConstructor
public class CustomAccessDeniedHandler implements AccessDeniedHandler {

    /**
     * ObjectMapper instance for converting Java objects to JSON format.
     */
    private final ObjectMapper objectMapper;

    /**
     * Handles the case where a user is denied access to a resource.
     * <p>
     * This method sets the response status to 403 (Forbidden), logs the access denial, and sends a JSON
     * error response with details about the exception.
     * </p>
     *
     * @param request               the {@link HttpServletRequest} object that triggered the exception
     * @param response              the {@link HttpServletResponse} object to which the error response is written
     * @param accessDeniedException the exception thrown when access is denied
     *
     * @throws IOException if an input or output exception occurs during response writing
     */
    @Override
    public void handle(HttpServletRequest request, HttpServletResponse response,
                       AccessDeniedException accessDeniedException) throws IOException {

        // Log the access denied exception with the relevant request information
        log.warn("Access denied for request to URL: {} | Reason: {}",
                 request.getRequestURI(),
                 accessDeniedException.getMessage());

        // Set response content type to JSON and the status code to 403 Forbidden
        response.setContentType("application/json");
        response.setStatus(HttpServletResponse.SC_FORBIDDEN);

        // Create error response body using a utility method
        var errorResponse = ResponseDtoUtils.createErrorResponseBody(null, HttpStatus.FORBIDDEN, accessDeniedException);

        // Write the JSON error response to the output stream
        response.getOutputStream().println(objectMapper.writeValueAsString(errorResponse));

        log.info("Forbidden access response sent to client for URL: {}", request.getRequestURI());
    }
}
