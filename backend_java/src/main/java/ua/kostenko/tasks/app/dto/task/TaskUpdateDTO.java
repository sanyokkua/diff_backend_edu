package ua.kostenko.tasks.app.dto.task;

import io.swagger.v3.oas.annotations.media.Schema;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * Data Transfer Object for updating a Task.
 * This class is used to encapsulate the data required to update an existing task, including the new name and description.
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
public class TaskUpdateDTO {

    /**
     * The new name of the task.
     */
    @Schema(description = "The new name of the task.", example = "Revise User Authentication Module")
    private String name;

    /**
     * The new description of the task.
     */
    @Schema(description = "The new description of the task.", example = "This task involves revising the user authentication module to improve security and user experience.")
    private String description;
}
