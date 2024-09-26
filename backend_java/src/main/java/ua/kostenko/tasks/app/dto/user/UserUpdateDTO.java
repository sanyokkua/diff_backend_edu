package ua.kostenko.tasks.app.dto.user;

import io.swagger.v3.oas.annotations.media.Schema;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * Data Transfer Object for updating user information.
 * This class is used to transfer user update data between processes.
 * It includes the user's current password, new password, and new password confirmation.
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
public class UserUpdateDTO {

    /**
     * The current password of the user.
     */
    @Schema(description = "The current password of the user.", example = "P@ssw0rd")
    private String currentPassword;

    /**
     * The new password of the user.
     */
    @Schema(description = "The new password of the user.", example = "N3wP@ssw0rd")
    private String newPassword;

    /**
     * The confirmation of the new password.
     */
    @Schema(description = "The confirmation of the new password.", example = "N3wP@ssw0rd")
    private String newPasswordConfirmation;
}
