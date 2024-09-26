package ua.kostenko.tasks.app.dto;

import com.fasterxml.jackson.annotation.JsonInclude;
import io.swagger.v3.oas.annotations.media.Schema;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import ua.kostenko.tasks.app.dto.task.TaskDto;
import ua.kostenko.tasks.app.dto.user.UserDto;

/**
 * A generic Data Transfer Object (DTO) used to wrap responses sent by the application.
 * <p>
 * This DTO contains status information, any response data, and potential error messages.
 * It is a generic class to allow flexible data types for the response payload.
 * </p>
 *
 * @param <T> The type of the data object included in the response.
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@JsonInclude(JsonInclude.Include.NON_NULL)
public class ResponseDto<T> {

    /**
     * The HTTP status code of the response.
     */
    @Schema(description = "The HTTP status code of the response.", example = "200")
    private int statusCode;

    /**
     * The HTTP status message corresponding to the status code.
     */
    @Schema(description = "The HTTP status message corresponding to the status code.", example = "OK")
    private String statusMessage;

    /**
     * The actual response data of type {@code T}.
     */
    @Schema(description = "The actual response data of type T.", implementation = Object.class, anyOf = {UserDto.class,
                                                                                                         TaskDto.class,
                                                                                                         String.class,
                                                                                                         Object.class})
    private T data;

    /**
     * Error message, if any, returned in case of an error response.
     * <p>
     * This field is only populated when an error occurs.
     * </p>
     */
    @Schema(description = "Error message, if any, returned in case of an error response.", example = "User not found.")
    private String error;
}
