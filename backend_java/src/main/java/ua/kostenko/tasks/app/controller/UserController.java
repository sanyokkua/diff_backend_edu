package ua.kostenko.tasks.app.controller;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.authentication.AuthenticationCredentialsNotFoundException;
import org.springframework.security.authentication.InsufficientAuthenticationException;
import org.springframework.web.bind.annotation.*;
import ua.kostenko.tasks.app.dto.ResponseDto;
import ua.kostenko.tasks.app.dto.user.UserDeletionDTO;
import ua.kostenko.tasks.app.dto.user.UserDto;
import ua.kostenko.tasks.app.dto.user.UserUpdateDTO;
import ua.kostenko.tasks.app.service.AuthUserExtractionService;
import ua.kostenko.tasks.app.service.UserService;
import ua.kostenko.tasks.app.utility.ResponseDtoUtils;
import ua.kostenko.tasks.app.utility.UserUtils;

/**
 * REST controller for managing user-related operations such as fetching user details,
 * updating passwords, and deleting users.
 *
 * <p>This controller provides endpoints for authenticated users to perform actions related
 * to their user data. It uses {@code AuthUserExtractionService} to extract the current
 * authenticated user from the security context, and {@code UserService} to handle user data operations.</p>
 *
 * <p>Each endpoint in this controller requires the user to be authenticated, otherwise,
 * appropriate exceptions such as {@link AuthenticationCredentialsNotFoundException} or
 * {@link InsufficientAuthenticationException} are thrown, preventing unauthorized access.</p>
 *
 * <p>Logging is used to capture the flow of execution, especially in case of errors, making
 * it easier to trace any issues during API usage.</p>
 */
@Tag(name = "User Management REST Controller", description = "Handles operations related to user data.")
@Slf4j
@RestController
@RequestMapping("api/v1/users")
@RequiredArgsConstructor
public class UserController {

    private final UserService userService;
    private final AuthUserExtractionService userExtractionService;

    /**
     * Retrieves the details of an authenticated user by ID.
     *
     * <p>This endpoint allows the authenticated user to retrieve their own details.
     * It uses the user ID from the path and compares it to the authenticated user's ID
     * to ensure that users can only fetch their own details.</p>
     *
     * @param userId the ID of the user to fetch
     *
     * @return the {@link ResponseEntity} containing the user details
     */
    @GetMapping("/{userId}")
    @Operation(summary = "Fetch User Details", description = "Retrieves the details of the authenticated user using their ID.")
    public ResponseEntity<ResponseDto<UserDto>> getUserById(
            @Parameter(description = "The ID of the user to fetch.") @PathVariable Long userId) {

        log.debug("Fetching user by ID: {}", userId);
        UserDto authenticatedUser = userExtractionService.getUserFromAuthContext();
        UserUtils.validateAuthenticatedUserIdWithPassed(authenticatedUser, userId);

        log.info("User fetched successfully: {}", authenticatedUser.getEmail());
        return ResponseDtoUtils.buildDtoResponse(authenticatedUser, HttpStatus.OK);
    }

    /**
     * Updates the password of the authenticated user.
     *
     * <p>This endpoint allows the user to update their password after providing the current
     * password and a new password. The current password is validated to ensure the user
     * can only change their password with proper authentication.</p>
     *
     * @param userId        the ID of the user whose password is being updated
     * @param userUpdateDTO the DTO containing the current and new passwords
     *
     * @return the {@link ResponseEntity} containing the updated user details
     */
    @PutMapping("/{userId}/password")
    @Operation(summary = "Update User Password", description = "Updates the password for the authenticated user.")
    public ResponseEntity<ResponseDto<UserDto>> updateUserPassword(
            @Parameter(description = "The ID of the user whose password is being updated.") @PathVariable Long userId,
            @Parameter(description = "DTO containing the current and new passwords.") @RequestBody UserUpdateDTO userUpdateDTO) {

        log.debug("Updating password for user ID: {}", userId);
        UserDto authenticatedUser = userExtractionService.getUserFromAuthContext();
        UserUtils.validateAuthenticatedUserIdWithPassed(authenticatedUser, userId);

        UserDto updatedUserDto = userService.updateUserPassword(userId, userUpdateDTO);
        log.info("Password updated successfully for user ID: {}", userId);
        return ResponseDtoUtils.buildDtoResponse(updatedUserDto, HttpStatus.OK);
    }

    /**
     * Deletes the authenticated user from the system.
     *
     * <p>This endpoint allows the authenticated user to delete their account by providing
     * the necessary credentials. The current password is validated before the deletion is allowed.</p>
     *
     * @param userId          the ID of the user to be deleted
     * @param userDeletionDTO the DTO containing deletion-related data
     *
     * @return the {@link ResponseEntity} confirming the deletion
     */
    @PostMapping("/{userId}/delete")
    @Operation(summary = "Delete User Account", description = "Deletes the authenticated user account from the system.")
    public ResponseEntity<ResponseDto<Void>> deleteUser(
            @Parameter(description = "The ID of the user to be deleted.") @PathVariable Long userId,
            @Parameter(description = "DTO containing deletion-related information.") @RequestBody UserDeletionDTO userDeletionDTO) {

        log.debug("Deleting user with ID: {}", userId);
        UserDto authenticatedUser = userExtractionService.getUserFromAuthContext();
        UserUtils.validateAuthenticatedUserIdWithPassed(authenticatedUser, userId);

        userService.deleteUser(userId, userDeletionDTO);
        log.info("User deleted successfully with ID: {}", userId);
        return ResponseDtoUtils.buildDtoResponse(null, HttpStatus.NO_CONTENT);
    }
}
