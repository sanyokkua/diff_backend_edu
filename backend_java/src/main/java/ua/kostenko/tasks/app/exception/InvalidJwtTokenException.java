package ua.kostenko.tasks.app.exception;

/**
 * Exception thrown when a JWT token is found to be invalid.
 * <p>
 * This exception is typically used in authentication processes to indicate that the provided JWT token
 * is either malformed, expired, or otherwise invalid.
 * </p>
 */
public class InvalidJwtTokenException extends RuntimeException {

    /**
     * Constructs a new {@code InvalidJwtTokenException} with the specified detail message.
     *
     * @param message the detail message, which is saved for later retrieval by the {@link #getMessage()} method.
     */
    public InvalidJwtTokenException(String message) {
        super(message);
    }
}
