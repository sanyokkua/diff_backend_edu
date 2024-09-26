package ua.kostenko.tasks.app.exception;

/**
 * Exception thrown when an attempt is made to create a task that already exists.
 * <p>
 * This exception is typically used in task management systems to indicate that the task being created
 * is a duplicate of an existing task.
 * </p>
 */
public class TaskAlreadyExistsException extends RuntimeException {

    /**
     * Constructs a new {@code TaskAlreadyExistsException} with the specified detail message.
     *
     * @param message the detail message, which is saved for later retrieval by the {@link #getMessage()} method.
     */
    public TaskAlreadyExistsException(String message) {
        super(message);
    }
}
