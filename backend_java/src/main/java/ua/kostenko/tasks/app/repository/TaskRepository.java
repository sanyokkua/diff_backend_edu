package ua.kostenko.tasks.app.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import ua.kostenko.tasks.app.entity.Task;
import ua.kostenko.tasks.app.entity.User;

import java.util.List;
import java.util.Optional;

/**
 * Repository interface for managing {@link Task} entities.
 * <p>
 * This interface extends {@link JpaRepository} to provide standard CRUD operations
 * for the Task entity and additional custom query methods.
 * </p>
 */
@Repository
public interface TaskRepository extends JpaRepository<Task, Long> {

    /**
     * Retrieves a task by the user and task name.
     *
     * @param user the user who owns the task
     * @param name the name of the task to be retrieved
     *
     * @return an Optional containing the Task if found, otherwise an empty Optional
     */
    Optional<Task> findByUserAndName(User user, String name);

    /**
     * Retrieves a task by the user and task ID.
     *
     * @param user   the user who owns the task
     * @param taskId the ID of the task to be retrieved
     *
     * @return an Optional containing the Task if found, otherwise an empty Optional
     */
    Optional<Task> findByUserAndTaskId(User user, Long taskId);

    /**
     * Retrieves all tasks associated with a specific user.
     *
     * @param user the user whose tasks are to be retrieved
     *
     * @return a list of tasks belonging to the user
     */
    List<Task> findAllByUser(User user);
}
