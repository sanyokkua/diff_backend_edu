package ua.kostenko.tasks.app.dto.user;

import io.swagger.v3.oas.annotations.media.Schema;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * Data Transfer Object for user information.
 * This class is used to transfer user data between processes.
 * It includes the user's ID, email, and JWT token.
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
public class UserDto {

    /**
     * The unique identifier of the user.
     */
    @Schema(description = "The unique identifier of the user.", example = "1")
    private Long userId;

    /**
     * The email of the user.
     */
    @Schema(description = "The email of the user.", example = "user@example.com")
    private String email;

    /**
     * The JWT token of the user.
     */
    @Schema(description = "The JWT token of the user.", example = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")
    private String jwtToken;
}
