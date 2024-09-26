package ua.kostenko.tasks.app.exception;

/**
 * Exception thrown when a password does not meet the required criteria.
 * <p>
 * This exception is typically used in validation processes to indicate that the provided password
 * does not comply with the security policies or format requirements.
 * </p>
 */
public class InvalidPasswordException extends RuntimeException {

    /**
     * Constructs a new {@code InvalidPasswordException} with the specified detail message.
     *
     * @param message the detail message, which is saved for later retrieval by the {@link #getMessage()} method.
     */
    public InvalidPasswordException(String message) {
        super(message);
    }
}
