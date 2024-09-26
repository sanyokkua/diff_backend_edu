package ua.kostenko.tasks.app.utility;

import lombok.AccessLevel;
import lombok.NoArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseCookie;
import org.springframework.http.ResponseEntity;
import ua.kostenko.tasks.app.dto.ResponseDto;

/**
 * Utility class for building standardized responses for the application.
 * <p>
 * This class provides methods to create response entities using {@link ResponseDto}
 * that encapsulate status codes, messages, data, and errors. It also supports building
 * responses with cookies and logging error details.
 * </p>
 */
@Slf4j
@NoArgsConstructor(access = AccessLevel.PRIVATE)
public class ResponseDtoUtils {

    /**
     * Constructs an error message from the exception.
     *
     * @param ex the exception to extract the error message from
     *
     * @return a formatted error message with the exception type and message
     */
    private static String getErrorMessage(Exception ex) {
        var type = ex.getClass().getSimpleName();
        var message = ex.getMessage();
        log.error("Error occurred: {} - {}", type, message); // Log the error
        return String.format("%s: %s", type, message);
    }

    /**
     * Builds a standard response DTO without any errors.
     * <p>
     * This method constructs a response entity containing the response data and status.
     * </p>
     *
     * @param <T>    the type of the response data
     * @param data   the response data to be included
     * @param status the HTTP status to be set in the response
     *
     * @return a {@link ResponseEntity} containing the {@link ResponseDto}
     */
    public static <T> ResponseEntity<ResponseDto<T>> buildDtoResponse(T data, HttpStatus status) {
        log.info("Building response with status: {} and data: {}", status, data); // Log status and data (debug mode)
        var responseBody =
                ResponseDto.<T>builder().data(data).statusCode(status.value()).statusMessage(status.name()).build();
        return ResponseEntity.status(status).body(responseBody);
    }

    /**
     * Builds a standard response DTO with a cookie.
     * <p>
     * This method constructs a response entity containing the response data, status,
     * and a cookie.
     * </p>
     *
     * @param <T>    the type of the response data
     * @param data   the response data to be included
     * @param status the HTTP status to be set in the response
     * @param cookie the {@link ResponseCookie} to be added in the response header
     *
     * @return a {@link ResponseEntity} containing the {@link ResponseDto}
     */
    public static <T> ResponseEntity<ResponseDto<T>> buildDtoResponse(T data, HttpStatus status,
                                                                      ResponseCookie cookie) {
        log.info("Building response with status: {}, data: {}, and cookie: {}",
                 status,
                 data,
                 cookie); // Log with cookie
        var responseBody =
                ResponseDto.<T>builder().data(data).statusCode(status.value()).statusMessage(status.name()).build();
        String cookieString = cookie.toString();
        return ResponseEntity.status(status).header(HttpHeaders.SET_COOKIE, cookieString).body(responseBody);
    }

    /**
     * Builds an error response DTO in case of an exception.
     * <p>
     * This method constructs a response entity that includes the exception's error message,
     * along with the response data and HTTP status.
     * </p>
     *
     * @param <T>    the type of the response data
     * @param data   the response data to be included
     * @param status the HTTP status to be set in the response
     * @param ex     the exception that occurred
     *
     * @return a {@link ResponseEntity} containing the error {@link ResponseDto}
     */
    public static <T> ResponseEntity<ResponseDto<T>> buildDtoErrorResponse(T data, HttpStatus status, Exception ex) {
        log.warn("Building error response due to exception: {}", ex.getMessage(), ex); // Log error details
        var responseBody = createErrorResponseBody(data, status, ex);
        return ResponseEntity.status(status).body(responseBody);
    }

    /**
     * Creates an error response body DTO in case of an exception.
     * <p>
     * This method builds a {@link ResponseDto} that includes the error message and status
     * based on the provided exception and response status.
     * </p>
     *
     * @param <T>    the type of the response data
     * @param data   the response data to be included
     * @param status the HTTP status to be set in the response
     * @param ex     the exception that occurred
     *
     * @return a {@link ResponseDto} containing the error message
     */
    public static <T> ResponseDto<T> createErrorResponseBody(T data, HttpStatus status, Exception ex) {
        var msg = getErrorMessage(ex);
        log.debug("Creating error response body with message: {}", msg); // Log error message for debug purposes
        return ResponseDto.<T>builder()
                          .data(data)
                          .statusCode(status.value())
                          .statusMessage(status.name())
                          .error(msg)
                          .build();
    }
}
