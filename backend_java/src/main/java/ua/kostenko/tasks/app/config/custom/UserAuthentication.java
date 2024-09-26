package ua.kostenko.tasks.app.config.custom;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.lang.Nullable;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.GrantedAuthority;
import ua.kostenko.tasks.app.entity.User;

import java.util.Collection;
import java.util.List;
import java.util.Objects;

/**
 * The {@code UserAuthentication} class is an implementation of the {@link Authentication} interface.
 * It represents the authentication information for a user, containing the user's details and JWT token.
 * <p>
 * This class uses the following Lombok annotations:
 * <ul>
 *   <li>{@code @Data} - Automatically generates getters, setters, toString, equals, and hashCode methods.</li>
 *   <li>{@code @Builder} - Provides a builder pattern for constructing instances of the class.</li>
 *   <li>{@code @RequiredArgsConstructor} - Generates a constructor for all final and non-initialized fields.</li>
 * </ul>
 * <p>
 * Key functionalities:
 * <ul>
 *   <li>Stores the user and their JWT token.</li>
 *   <li>Implements the core methods of the {@link Authentication} interface for handling user credentials and authorities.</li>
 *   <li>Ensures that passwords are not stored or exposed, returning {@code null} for credentials.</li>
 * </ul>
 * <p>
 * Note: {@code transient} is used for fields that shouldn't be serialized (e.g., during sessions).
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class UserAuthentication implements Authentication {

    /**
     * The {@code User} object representing the authenticated user.
     * Marked as {@code transient} to avoid serialization during session handling.
     */
    private transient User user;

    /**
     * The JWT token associated with this authentication.
     * Also marked as {@code transient} to avoid serialization.
     */
    private transient String jwtToken;

    /**
     * Boolean flag to indicate whether the user is authenticated.
     */
    private boolean isAuthenticated;

    /**
     * {@inheritDoc}
     * <p>
     * Returns the name of the authenticated user. The name is represented by the user's email address.
     * If the user is not set, an empty string is returned.
     * </p>
     *
     * @return The email of the authenticated user, or an empty string if no user is authenticated.
     */
    @Override
    public String getName() {
        return Objects.nonNull(user) ? user.getEmail() : "";
    }

    /**
     * {@inheritDoc}
     * <p>
     * Returns the authorities granted to the authenticated user.
     * This implementation returns an empty list, indicating no specific authorities.
     * </p>
     *
     * @return An empty collection, since this implementation does not use authorities.
     */
    @Override
    public Collection<? extends GrantedAuthority> getAuthorities() {
        return List.of();
    }

    /**
     * {@inheritDoc}
     * <p>
     * This method returns {@code null} since the application doesn't use or store passwords.
     * Password-based authentication is not needed.
     * </p>
     *
     * @return Always returns {@code null}.
     */
    @Nullable
    @Override
    public String getCredentials() {
        // Password-based authentication is not used, hence null is returned.
        return null;
    }

    /**
     * {@inheritDoc}
     * <p>
     * Returns additional details about the authentication, in this case, the JWT token.
     * If the JWT token is {@code null}, an empty string is returned.
     * </p>
     *
     * @return The JWT token, or an empty string if the token is not available.
     */
    @Override
    public String getDetails() {
        return Objects.nonNull(jwtToken) ? jwtToken : "";
    }

    /**
     * {@inheritDoc}
     * <p>
     * This method returns the {@code User} object representing the principal (authenticated user).
     * </p>
     *
     * @return The authenticated {@code User}, or {@code null} if no user is set.
     */
    @Nullable
    @Override
    public User getPrincipal() {
        return user;
    }

    /**
     * {@inheritDoc}
     * <p>
     * Checks whether the user is authenticated.
     * </p>
     *
     * @return {@code true} if the user is authenticated, {@code false} otherwise.
     */
    @Override
    public boolean isAuthenticated() {
        return isAuthenticated;
    }

    /**
     * {@inheritDoc}
     * <p>
     * Sets the authentication status of the user. Throws an {@link IllegalArgumentException}
     * if attempting to set an invalid value (such as {@code true} for an anonymous or unauthenticated user).
     * </p>
     *
     * @param isAuthenticated the new authentication status.
     *
     * @throws IllegalArgumentException if the input is invalid.
     */
    @Override
    public void setAuthenticated(boolean isAuthenticated) throws IllegalArgumentException {
        this.isAuthenticated = isAuthenticated;
    }
}
