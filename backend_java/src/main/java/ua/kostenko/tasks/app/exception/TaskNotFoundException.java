package ua.kostenko.tasks.app.exception;

/**
 * Exception thrown when a requested task cannot be found.
 * <p>
 * This exception is typically used in task management systems to indicate that the task being queried
 * does not exist in the system.
 * </p>
 */
public class TaskNotFoundException extends RuntimeException {

    /**
     * Constructs a new {@code TaskNotFoundException} with the specified detail message.
     *
     * @param message the detail message, which is saved for later retrieval by the {@link #getMessage()} method.
     */
    public TaskNotFoundException(String message) {
        super(message);
    }
}
