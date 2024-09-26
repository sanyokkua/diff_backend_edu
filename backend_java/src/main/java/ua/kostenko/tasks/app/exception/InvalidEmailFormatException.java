package ua.kostenko.tasks.app.exception;

/**
 * Exception thrown when an email address does not conform to the expected format.
 * <p>
 * This exception is typically used in validation processes to indicate that the provided email
 * address does not match the required format.
 * </p>
 */
public class InvalidEmailFormatException extends RuntimeException {

    /**
     * Constructs a new {@code InvalidEmailFormatException} with the specified detail message.
     *
     * @param message the detail message, which is saved for later retrieval by the {@link #getMessage()} method.
     */
    public InvalidEmailFormatException(String message) {
        super(message);
    }
}
