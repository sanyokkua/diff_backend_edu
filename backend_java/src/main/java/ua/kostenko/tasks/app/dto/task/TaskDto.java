package ua.kostenko.tasks.app.dto.task;

import io.swagger.v3.oas.annotations.media.Schema;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * Data Transfer Object for Task.
 * This class is used to encapsulate the data related to a task, including its ID, name, description,
 * and the ID of the user associated with the task.
 *
 * <p>
 * Annotations used:
 * <ul>
 *   <li>{@link Data} - Generates getters, setters, toString, equals, and hashCode methods.</li>
 *   <li>{@link Builder} - Implements the builder pattern for object creation.</li>
 *   <li>{@link NoArgsConstructor} - Generates a no-argument constructor.</li>
 *   <li>{@link AllArgsConstructor} - Generates a constructor with one parameter for each field in the class.</li>
 * </ul>
 * </p>
 *
 * @see lombok.Data
 * @see lombok.Builder
 * @see lombok.NoArgsConstructor
 * @see lombok.AllArgsConstructor
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class TaskDto {

    /**
     * The unique identifier of the task.
     */
    @Schema(description = "The unique identifier of the task.", example = "1")
    private Long taskId;

    /**
     * The name of the task.
     */
    @Schema(description = "The name of the task.", example = "Implement OpenAPI Documentation")
    private String name;

    /**
     * A brief description of the task.
     */
    @Schema(description = "A brief description of the task.", example = "This task involves creating OpenAPI documentation for the Spring Boot application.")
    private String description;

    /**
     * The unique identifier of the user associated with the task.
     */
    @Schema(description = "The unique identifier of the user associated with the task.", example = "123")
    private Long userId;
}
