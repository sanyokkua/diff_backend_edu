package ua.kostenko.tasks.app.config.custom;

import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.web.AuthenticationEntryPoint;
import org.springframework.stereotype.Component;
import ua.kostenko.tasks.app.utility.ResponseDtoUtils;

import java.io.IOException;

/**
 * Custom Authentication Entry Point for handling unauthorized access in Spring Security.
 * <p>
 * This component is invoked whenever an unauthenticated user attempts to access a secured resource.
 * It sends a JSON response with a 401 Unauthorized status code, along with details of the authentication
 * failure.
 * </p>
 */
@Slf4j
@Component
@RequiredArgsConstructor
public class CustomAuthenticationEntryPoint implements AuthenticationEntryPoint {

    /**
     * ObjectMapper instance for converting Java objects to JSON format.
     */
    private final ObjectMapper objectMapper;

    /**
     * Commences an authentication scheme.
     * <p>
     * This method is called when an unauthenticated user tries to access a secured endpoint.
     * It sets the response status to 401 (Unauthorized) and writes a JSON error response.
     * </p>
     *
     * @param request       the {@link HttpServletRequest} object that triggered the authentication exception
     * @param response      the {@link HttpServletResponse} object to which the error response is written
     * @param authException the exception that triggered the authentication process
     *
     * @throws IOException if an input or output exception occurs during response writing
     */
    @Override
    public void commence(HttpServletRequest request, HttpServletResponse response,
                         AuthenticationException authException) throws IOException {

        // Log the unauthorized access attempt with details of the request
        log.warn("Unauthorized access attempt to URL: {} | Reason: {}",
                 request.getRequestURI(),
                 authException.getMessage());

        // Set the response content type to JSON and the status code to 401 Unauthorized
        response.setContentType("application/json");
        response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);

        // Create error response body using a utility method
        var errorResponse = ResponseDtoUtils.createErrorResponseBody(null, HttpStatus.UNAUTHORIZED, authException);

        // Write the JSON error response to the output stream
        response.getOutputStream().println(objectMapper.writeValueAsString(errorResponse));

        log.info("Unauthorized access response sent to client for URL: {}", request.getRequestURI());
    }
}
