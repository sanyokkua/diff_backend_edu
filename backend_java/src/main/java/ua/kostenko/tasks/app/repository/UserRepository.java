package ua.kostenko.tasks.app.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import ua.kostenko.tasks.app.entity.User;

import java.util.Optional;

/**
 * Repository interface for managing {@link User} entities.
 * <p>
 * This interface extends {@link JpaRepository} to provide standard CRUD operations
 * for the User entity and additional custom query methods.
 * </p>
 */
@Repository
public interface UserRepository extends JpaRepository<User, Long> {

    /**
     * Retrieves a user by their email address.
     *
     * @param email the email address of the user to be retrieved
     *
     * @return an Optional containing the User if found, otherwise an empty Optional
     */
    Optional<User> findByEmail(String email);
}
