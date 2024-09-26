package ua.kostenko.tasks.app.dto.user;

import io.swagger.v3.oas.annotations.media.Schema;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * Data Transfer Object for user login.
 * This class is used to transfer user login data between processes.
 * It includes the user's email and password.
 *
 * <p>
 * Annotations used:
 * <ul>
 *   <li>{@code @Data} - Generates getters, setters, toString, equals, and hashCode methods.</li>
 *   <li>{@code @Builder} - Implements the builder pattern for object creation.</li>
 *   <li>{@code @NoArgsConstructor} - Generates a no-argument constructor.</li>
 *   <li>{@code @AllArgsConstructor} - Generates a constructor with all fields as parameters.</li>
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
public class UserLoginDto {

    /**
     * The email of the user.
     */
    @Schema(description = "The email of the user.", example = "user@example.com")
    private String email;

    /**
     * The password of the user.
     */
    @Schema(description = "The password of the user.", example = "P@ssw0rd")
    private String password;
}
