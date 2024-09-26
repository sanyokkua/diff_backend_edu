package ua.kostenko.tasks.app.entity;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.HashSet;
import java.util.Set;

/**
 * Represents a user entity in the system.
 * <p>
 * This class maps to the "users" table in the database and contains information
 * about the user's email, password hash, and associated tasks.
 * </p>
 *
 * <p>
 * A user must have a unique email address, and the password is stored as a hash
 * for security reasons. The relationship with tasks is one-to-many, where one user
 * can have multiple tasks.
 * </p>
 */
@Entity
@Table(name = "users", schema = "backend_diff", uniqueConstraints = {@UniqueConstraint(columnNames = "email")})
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class User {

    /**
     * The unique identifier for the user.
     * <p>
     * This value is automatically generated by the database upon user creation.
     * </p>
     */
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long userId;

    /**
     * The user's email address.
     * <p>
     * This field is mandatory and must be unique across all users.
     * </p>
     */
    @Column(nullable = false, unique = true)
    private String email;

    /**
     * The hashed password of the user.
     * <p>
     * This field is mandatory and is used for user authentication.
     * </p>
     */
    @Column(nullable = false)
    private String passwordHash;

    /**
     * The set of tasks associated with the user.
     * <p>
     * This field establishes a one-to-many relationship with the Task entity.
     * When a user is deleted, all associated tasks are also removed.
     * </p>
     */
    @OneToMany(mappedBy = "user", cascade = CascadeType.ALL, orphanRemoval = true)
    private Set<Task> tasks = new HashSet<>();
}
