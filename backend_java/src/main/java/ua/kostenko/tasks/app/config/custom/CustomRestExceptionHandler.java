package ua.kostenko.tasks.app.config.custom;

import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.AccessDeniedException;
import org.springframework.security.authentication.AuthenticationCredentialsNotFoundException;
import org.springframework.security.authentication.InsufficientAuthenticationException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;
import org.springframework.web.servlet.NoHandlerFoundException;
import ua.kostenko.tasks.app.dto.ResponseDto;
import ua.kostenko.tasks.app.exception.*;
import ua.kostenko.tasks.app.utility.ResponseDtoUtils;

import java.io.IOException;

/**
 * Global exception handler for REST APIs using @RestControllerAdvice to catch and handle all application-level exceptions.
 * This class is responsible for converting various exceptions into appropriate HTTP responses with error details.
 */
@SuppressWarnings("GrazieInspection")
@Slf4j
@RestControllerAdvice
public class CustomRestExceptionHandler {

    /**
     * Builds a standardized error response including the exception message, status, and request details.
     *
     * @param ex     the exception thrown
     * @param req    the HttpServletRequest where the exception occurred
     * @param status the corresponding HTTP status code
     *
     * @return ResponseEntity containing the error details wrapped in a ResponseDto object
     */
    private ResponseEntity<ResponseDto<Object>> buildErrorResponse(Exception ex, HttpServletRequest req,
                                                                   HttpStatus status) {
        log.warn("Handling exception: {} | Status: {} | Path: {}",
                 ex.getClass().getSimpleName(),
                 status,
                 req.getRequestURI(),
                 ex);
        var requestBody = req.getAttribute("requestBody");  // Custom requestBody extraction if available
        return ResponseDtoUtils.buildDtoErrorResponse(requestBody, status, ex);
    }

    /**
     * Handles all exceptions that indicate invalid requests from the client.
     * Maps to HttpStatus.BAD_REQUEST (400).
     *
     * @param ex      the caught exception
     * @param request the incoming HttpServletRequest
     *
     * @return ResponseEntity with a 400 Bad Request status and error details
     */
    @ExceptionHandler({IllegalArgumentException.class,
                       InvalidEmailFormatException.class,
                       InvalidPasswordException.class,
                       InvalidJwtTokenException.class})
    public ResponseEntity<ResponseDto<Object>> handleBadRequestException(Exception ex, HttpServletRequest request) {
        return buildErrorResponse(ex, request, HttpStatus.BAD_REQUEST);
    }

    /**
     * Handles exceptions indicating an unauthorized request, such as missing or invalid authentication credentials.
     * Maps to HttpStatus.UNAUTHORIZED (401).
     *
     * @param ex      the caught exception
     * @param request the incoming HttpServletRequest
     *
     * @return ResponseEntity with a 401 Unauthorized status and error details
     */
    @ExceptionHandler({AuthenticationCredentialsNotFoundException.class, InsufficientAuthenticationException.class})
    public ResponseEntity<ResponseDto<Object>> handleUnauthorizedException(Exception ex, HttpServletRequest request) {
        return buildErrorResponse(ex, request, HttpStatus.UNAUTHORIZED);
    }

    /**
     * Handles access denied exceptions for cases where the client has insufficient permissions.
     * Maps to HttpStatus.FORBIDDEN (403).
     *
     * @param ex      the caught exception
     * @param request the incoming HttpServletRequest
     *
     * @return ResponseEntity with a 403 Forbidden status and error details
     */
    @ExceptionHandler({AccessDeniedException.class})
    public ResponseEntity<ResponseDto<Object>> handleAccessDeniedException(Exception ex, HttpServletRequest request) {
        return buildErrorResponse(ex, request, HttpStatus.FORBIDDEN);
    }

    /**
     * Handles exceptions when requested resources are not found.
     * Maps to HttpStatus.NOT_FOUND (404).
     *
     * @param ex      the caught exception
     * @param request the incoming HttpServletRequest
     *
     * @return ResponseEntity with a 404 Not Found status and error details
     */
    @ExceptionHandler({TaskNotFoundException.class})
    public ResponseEntity<ResponseDto<Object>> handleNotFoundException(Exception ex, HttpServletRequest request) {
        return buildErrorResponse(ex, request, HttpStatus.NOT_FOUND);
    }

    /**
     * Handles requests made to invalid endpoints or resources that don't exist.
     * Maps to HttpStatus.NOT_FOUND (404) for NoHandlerFoundException.
     *
     * @param ex      the NoHandlerFoundException exception
     * @param request the incoming HttpServletRequest
     *
     * @return ResponseEntity with a 404 Not Found status and error details
     */
    @ExceptionHandler(NoHandlerFoundException.class)
    public ResponseEntity<ResponseDto<Object>> handleNoHandlerFoundException(NoHandlerFoundException ex,
                                                                             HttpServletRequest request) {
        return buildErrorResponse(ex, request, HttpStatus.NOT_FOUND);
    }

    /**
     * Handles exceptions where resources (e.g., email, task) already exist, causing a conflict.
     * Maps to HttpStatus.CONFLICT (409).
     *
     * @param ex      the caught exception
     * @param request the incoming HttpServletRequest
     *
     * @return ResponseEntity with a 409 Conflict status and error details
     */
    @ExceptionHandler({EmailAlreadyExistsException.class, TaskAlreadyExistsException.class})
    public ResponseEntity<ResponseDto<Object>> handleConflictException(Exception ex, HttpServletRequest request) {
        return buildErrorResponse(ex, request, HttpStatus.CONFLICT);
    }

    /**
     * Handles general runtime exceptions or any uncaught exceptions that are not explicitly handled by other methods.
     * Maps to HttpStatus.INTERNAL_SERVER_ERROR (500).
     *
     * @param ex      the caught exception
     * @param request the incoming HttpServletRequest
     *
     * @return ResponseEntity with a 500 Internal Server Error status and error details
     */
    @ExceptionHandler(RuntimeException.class)
    public ResponseEntity<ResponseDto<Object>> handleRuntimeException(RuntimeException ex, HttpServletRequest request) {
        return buildErrorResponse(ex, request, HttpStatus.INTERNAL_SERVER_ERROR);
    }

    /**
     * Handles IO exceptions that can occur due to I/O issues.
     * Maps to HttpStatus.INTERNAL_SERVER_ERROR (500).
     *
     * @param ex      the IOException exception
     * @param request the incoming HttpServletRequest
     *
     * @return ResponseEntity with a 500 Internal Server Error status and error details
     */
    @ExceptionHandler(IOException.class)
    public ResponseEntity<ResponseDto<Object>> handleIOException(IOException ex, HttpServletRequest request) {
        return buildErrorResponse(ex, request, HttpStatus.INTERNAL_SERVER_ERROR);
    }

    /**
     * Handles ServletException for general servlet-related issues.
     * Maps to HttpStatus.INTERNAL_SERVER_ERROR (500).
     *
     * @param ex      the ServletException exception
     * @param request the incoming HttpServletRequest
     *
     * @return ResponseEntity with a 500 Internal Server Error status and error details
     */
    @ExceptionHandler(ServletException.class)
    public ResponseEntity<ResponseDto<Object>> handleServletException(ServletException ex, HttpServletRequest request) {
        return buildErrorResponse(ex, request, HttpStatus.INTERNAL_SERVER_ERROR);
    }

    /**
     * Catches any other exceptions not covered by specific handlers, acting as a global fallback.
     * Maps to HttpStatus.INTERNAL_SERVER_ERROR (500).
     *
     * @param ex      the caught exception
     * @param request the incoming HttpServletRequest
     *
     * @return ResponseEntity with a 500 Internal Server Error status and error details
     */
    @ExceptionHandler(Exception.class)
    public ResponseEntity<ResponseDto<Object>> handleGenericException(Exception ex, HttpServletRequest request) {
        return buildErrorResponse(ex, request, HttpStatus.INTERNAL_SERVER_ERROR);
    }
}
