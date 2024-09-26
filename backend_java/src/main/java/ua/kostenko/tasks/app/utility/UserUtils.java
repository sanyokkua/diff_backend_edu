package ua.kostenko.tasks.app.utility;

import lombok.AccessLevel;
import lombok.NoArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.access.AccessDeniedException;
import org.springframework.security.crypto.password.PasswordEncoder;
import ua.kostenko.tasks.app.dto.user.*;
import ua.kostenko.tasks.app.entity.User;
import ua.kostenko.tasks.app.exception.EmailAlreadyExistsException;
import ua.kostenko.tasks.app.exception.InvalidEmailFormatException;
import ua.kostenko.tasks.app.exception.InvalidPasswordException;
import ua.kostenko.tasks.app.repository.UserRepository;

import java.util.Objects;
import java.util.regex.Pattern;

/**
 * Utility class for user-related operations such as email validation, password validation, and existence checks.
 * <p>
 * This class provides static methods for various user-related utility functions such as:
 * <ul>
 *   <li>Validating email format</li>
 *   <li>Checking if a user already exists based on email</li>
 *   <li>Validating password matching and updates</li>
 *   <li>Checking authorization for user-related actions</li>
 * </ul>
 * </p>
 * <p>
 * It is marked as {@code @NoArgsConstructor(access = AccessLevel.PRIVATE)} to ensure it cannot be instantiated,
 * following the Singleton pattern for utility classes.
 * </p>
 * <p>
 * This class also uses the {@code @Slf4j} annotation from Lombok to provide logging capabilities.
 * </p>
 */
@Slf4j
@NoArgsConstructor(access = AccessLevel.PRIVATE)
public class UserUtils {

    // Regular expression pattern for validating email format
    private static final String EMAIL_PATTERN = "^[\\w-.]+@([\\w-]+\\.)[\\w-]{2,4}$";
    private static final Pattern pattern = Pattern.compile(EMAIL_PATTERN);

    /**
     * Validates the format of an email address.
     * <p>
     * This method checks if the provided email matches the standard email format using a regular expression.
     * If the format is invalid, it logs a warning and throws {@link InvalidEmailFormatException}.
     * </p>
     *
     * @param email the email to validate
     *
     * @throws InvalidEmailFormatException if the email format is invalid
     */
    public static void validateEmailFormat(String email) {
        log.debug("Validating email format for: {}", email);
        if (Objects.isNull(email) || !pattern.matcher(email).matches()) {
            log.warn("Invalid email format: {}", email);
            throw new InvalidEmailFormatException("Invalid email format");
        }
        log.info("Email format is valid: {}", email);
    }

    /**
     * Checks if a user with the given email exists in the repository.
     * <p>
     * This method queries the {@link UserRepository} to check if an email is already associated with a user.
     * If the email exists, it logs a warning and throws {@link EmailAlreadyExistsException}.
     * </p>
     *
     * @param userRepository the repository to query for user existence
     * @param email          the email to check
     *
     * @throws EmailAlreadyExistsException if the email already exists
     */
    public static void checkUserExists(UserRepository userRepository, String email) {
        log.debug("Checking if user exists with email: {}", email);
        userRepository.findByEmail(email).ifPresent(user -> {
            log.warn("Email already exists: {}", email);
            throw new EmailAlreadyExistsException("Email is already in use");
        });
        log.info("Email is available for use: {}", email);
    }

    /**
     * Validates if the provided password matches its confirmation.
     * <p>
     * This method checks if the password and the password confirmation are equal.
     * If they don't match, it logs a warning and throws {@link InvalidPasswordException}.
     * </p>
     *
     * @param password             the password to validate
     * @param passwordConfirmation the confirmation password
     *
     * @throws InvalidPasswordException if passwords do not match
     */
    public static void validatePasswords(String password, String passwordConfirmation) {
        log.debug("Validating password confirmation.");
        if (Objects.isNull(password) || Objects.isNull(passwordConfirmation) || !password.equals(passwordConfirmation)) {
            log.warn("Passwords do not match.");
            throw new InvalidPasswordException("Passwords do not match");
        }
        log.info("Passwords match.");
    }

    /**
     * Validates a password update request.
     * <p>
     * This method validates the current password, new password, and their confirmation.
     * It ensures that the current password is correct, the new password isn't the same as the current one,
     * and that the new password matches its confirmation. Logs are added for both success and failure cases.
     * </p>
     *
     * @param userUpdateDTO   the DTO containing current and new passwords
     * @param user            the user entity for which the password is being updated
     * @param passwordEncoder the password encoder used to verify the current password
     *
     * @throws InvalidPasswordException if any password validation fails
     */
    public static void validatePasswordUpdate(UserUpdateDTO userUpdateDTO, User user, PasswordEncoder passwordEncoder) {
        log.debug("Validating password update for user: {}", user.getEmail());

        if (!passwordEncoder.matches(userUpdateDTO.getCurrentPassword(), user.getPasswordHash())) {
            log.warn("Current password is incorrect for user: {}", user.getEmail());
            throw new InvalidPasswordException("Current password is incorrect");
        }

        if (userUpdateDTO.getNewPassword().equals(userUpdateDTO.getCurrentPassword())) {
            log.warn("New password cannot be the same as the current password for user: {}", user.getEmail());
            throw new InvalidPasswordException("New password cannot be the same as the current password");
        }

        validatePasswords(userUpdateDTO.getNewPassword(), userUpdateDTO.getNewPasswordConfirmation());
        log.info("Password update validation successful for user: {}", user.getEmail());
    }

    /**
     * Validates that the authenticated user's ID matches the ID passed in the request.
     * <p>
     * This is used to verify that a user is authorized to perform an action on a resource.
     * If the IDs do not match, an {@link AccessDeniedException} is thrown.
     * </p>
     *
     * @param userDto the authenticated user's data
     * @param userId  the user ID passed in the request
     *
     * @throws AccessDeniedException if the user is not authorized to perform the action
     */
    public static void validateAuthenticatedUserIdWithPassed(UserDto userDto, Long userId) {
        log.debug("Validating if authenticated user ID matches passed user ID.");
        if (Objects.isNull(userDto) || Objects.isNull(userId) || !Objects.equals(userDto.getUserId(), userId)) {
            log.warn("User ID validation failed: Authenticated user ID does not match passed user ID.");
            throw new AccessDeniedException("User is not authorized to perform this action");
        }
        log.info("User ID validation successful.");
    }

    /**
     * Validates the {@link UserLoginDto} for non-null and non-empty fields.
     * <p>
     * This method checks if the provided email and password in the {@code UserLoginDto} are valid.
     * Logs an error if validation fails.
     * </p>
     *
     * @param userLoginDto the login DTO to validate
     *
     * @throws IllegalArgumentException if the DTO or its fields are invalid
     */
    public static void validateUserLoginDto(UserLoginDto userLoginDto) {
        log.debug("Validating UserLoginDto.");
        if (Objects.isNull(userLoginDto)) {
            log.error("UserLoginDto is null.");
            throw new IllegalArgumentException("UserLoginDto is null");
        }
        if (Objects.isNull(userLoginDto.getEmail()) || userLoginDto.getEmail().isBlank()) {
            log.error("UserLoginDto email is null or empty.");
            throw new IllegalArgumentException("UserLoginDto email is null or empty");
        }
        if (Objects.isNull(userLoginDto.getPassword()) || userLoginDto.getPassword().isBlank()) {
            log.error("UserLoginDto password is null or empty.");
            throw new IllegalArgumentException("UserLoginDto password is null or empty");
        }
        log.info("UserLoginDto validation successful.");
    }

    /**
     * Validates the {@link UserCreationDTO} for non-null and non-empty fields.
     * <p>
     * This method ensures that all required fields (email, password, password confirmation) are present
     * and logs validation steps.
     * </p>
     *
     * @param userCreationDTO the user creation DTO to validate
     *
     * @throws IllegalArgumentException if any of the required fields are invalid
     */
    public static void validateUserCreationDTO(UserCreationDTO userCreationDTO) {
        log.debug("Validating UserCreationDTO.");
        if (Objects.isNull(userCreationDTO)) {
            log.error("UserCreationDTO is null.");
            throw new IllegalArgumentException("UserCreationDTO is null");
        }
        if (Objects.isNull(userCreationDTO.getEmail()) || userCreationDTO.getEmail().isBlank()) {
            log.error("UserCreationDTO email is null or empty.");
            throw new IllegalArgumentException("UserCreationDTO email is null or empty");
        }
        if (Objects.isNull(userCreationDTO.getPassword()) || userCreationDTO.getPassword().isBlank()) {
            log.error("UserCreationDTO password is null or empty.");
            throw new IllegalArgumentException("UserCreationDTO password is null or empty");
        }
        if (Objects.isNull(userCreationDTO.getPasswordConfirmation()) || userCreationDTO.getPasswordConfirmation()
                                                                                        .isBlank()) {
            log.error("UserCreationDTO password confirmation is null or empty.");
            throw new IllegalArgumentException("UserCreationDTO password confirmation is null or empty");
        }

        validatePasswords(userCreationDTO.getPassword(), userCreationDTO.getPasswordConfirmation());
        log.info("UserCreationDTO validation successful.");
    }

    /**
     * Validates the {@link UserUpdateDTO} for non-null and non-empty fields.
     * <p>
     * This method ensures that the current password, new password, and confirmation are valid.
     * Logs validation steps and failures.
     * </p>
     *
     * @param userUpdateDTO the user update DTO to validate
     *
     * @throws IllegalArgumentException if any required fields are invalid
     */
    public static void validateUserUpdateDTO(UserUpdateDTO userUpdateDTO) {
        log.debug("Validating UserUpdateDTO.");
        if (Objects.isNull(userUpdateDTO)) {
            log.error("UserUpdateDTO is null.");
            throw new IllegalArgumentException("UserUpdateDTO is null");
        }
        if (Objects.isNull(userUpdateDTO.getCurrentPassword()) || userUpdateDTO.getCurrentPassword().isBlank()) {
            log.error("UserUpdateDTO currentPassword is null or empty.");
            throw new IllegalArgumentException("UserUpdateDTO currentPassword is null or empty");
        }
        if (Objects.isNull(userUpdateDTO.getNewPassword()) || userUpdateDTO.getNewPassword().isBlank()) {
            log.error("UserUpdateDTO newPassword is null or empty.");
            throw new IllegalArgumentException("UserUpdateDTO newPassword is null or empty");
        }
        if (Objects.isNull(userUpdateDTO.getNewPasswordConfirmation()) || userUpdateDTO.getNewPasswordConfirmation()
                                                                                       .isBlank()) {
            log.error("UserUpdateDTO newPasswordConfirmation is null or empty.");
            throw new IllegalArgumentException("UserUpdateDTO newPasswordConfirmation is null or empty");
        }
        log.info("UserUpdateDTO validation successful.");
    }

    /**
     * Validates the {@link UserDeletionDTO} for non-null and non-empty fields.
     * <p>
     * This method ensures that the email and password in the user deletion DTO are valid.
     * Logs validation steps and failures.
     * </p>
     *
     * @param userDeletionDTO the user deletion DTO to validate
     *
     * @throws IllegalArgumentException if any required fields are invalid
     */
    public static void validateUserDeletionDTO(UserDeletionDTO userDeletionDTO) {
        log.debug("Validating UserDeletionDTO.");
        if (Objects.isNull(userDeletionDTO)) {
            log.error("UserDeletionDTO is null.");
            throw new IllegalArgumentException("UserDeletionDTO is null");
        }
        if (Objects.isNull(userDeletionDTO.getEmail()) || userDeletionDTO.getEmail().isBlank()) {
            log.error("UserDeletionDTO email is null or empty.");
            throw new IllegalArgumentException("UserDeletionDTO email is null or empty");
        }
        if (Objects.isNull(userDeletionDTO.getCurrentPassword()) || userDeletionDTO.getCurrentPassword().isBlank()) {
            log.error("UserDeletionDTO currentPassword is null or empty.");
            throw new IllegalArgumentException("UserDeletionDTO password is null or empty");
        }
        log.info("UserDeletionDTO validation successful.");
    }
}
