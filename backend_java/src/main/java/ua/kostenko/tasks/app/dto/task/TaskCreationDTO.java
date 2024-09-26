package ua.kostenko.tasks.app.dto.task;

import io.swagger.v3.oas.annotations.media.Schema;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * Data Transfer Object for creating a Task.
 * This class is used to encapsulate the data required to create a new task.
 * It includes the task's name and description.
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
public class TaskCreationDTO {

    /**
     * The name of the task.
     */
    @Schema(description = "The name of the task.", example = "Develop User Authentication Module")
    private String name;

    /**
     * A brief description of the task.
     */
    @Schema(description = "A brief description of the task.", example = "This task involves developing the user authentication module for the application, including login and registration functionality.")
    private String description;
}
