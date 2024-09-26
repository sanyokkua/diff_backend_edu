package ua.kostenko.tasks.app.exception;

/**
 * Exception thrown when an attempt is made to register an email that already exists in the system.
 * <p>
 * This exception is typically used in user registration processes to indicate that the provided email
 * address is already associated with an existing account.
 * </p>
 */
public class EmailAlreadyExistsException extends RuntimeException {

    /**
     * Constructs a new {@code EmailAlreadyExistsException} with the specified detail message.
     *
     * @param message the detail message, which is saved for later retrieval by the {@link #getMessage()} method.
     */
    public EmailAlreadyExistsException(String message) {
        super(message);
    }
}
