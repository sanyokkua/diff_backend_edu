package ua.kostenko.tasks.app.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.authentication.AuthenticationCredentialsNotFoundException;
import org.springframework.security.authentication.InsufficientAuthenticationException;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Service;
import ua.kostenko.tasks.app.config.custom.UserAuthentication;
import ua.kostenko.tasks.app.dto.user.UserDto;
import ua.kostenko.tasks.app.entity.User;

import java.util.Optional;

/**
 * The {@code AuthUserExtractionService} is a service class responsible for extracting
 * {@link UserDto} from the current security context's {@link Authentication} object.
 * <p>
 * This class ensures that the user is properly authenticated and that the required user
 * information is available. If the user is not authenticated or the required data is missing,
 * it throws appropriate exceptions.
 * </p>
 * <p>
 * This class uses Lombok annotations:
 * <ul>
 *   <li>{@code @Service} - Marks the class as a Spring service component.</li>
 *   <li>{@code @Slf4j} - Provides a logger instance for logging operations.</li>
 *   <li>{@code @RequiredArgsConstructor} - Generates a constructor for all final fields.</li>
 * </ul>
 * </p>
 */
@Service
@Slf4j
@RequiredArgsConstructor
public class AuthUserExtractionService {

    /**
     * Extracts the {@link UserDto} from the current authentication context.
     * <p>
     * This method retrieves the current {@link Authentication} object from the {@link SecurityContextHolder}.
     * If the authentication is not of type {@link UserAuthentication}, it logs a warning and throws an
     * {@link AuthenticationCredentialsNotFoundException}.
     * Also {@link InsufficientAuthenticationException} can be thrown during principal extraction.
     * </p>
     *
     * @return the {@link UserDto} extracted from the authentication context.
     *
     * @throws AuthenticationCredentialsNotFoundException if the user is not authenticated.
     * @throws InsufficientAuthenticationException        if the principal is not set. (extractUserDto())
     */
    public UserDto getUserFromAuthContext() {
        log.debug("Attempting to retrieve user from authentication context.");

        // Retrieve the current authentication object
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();

        // Check if authentication is of type UserAuthentication
        if (!(authentication instanceof UserAuthentication userAuthentication)) {
            log.warn("User is not authenticated. Authentication object: {}", authentication);
            throw new AuthenticationCredentialsNotFoundException("User is not authenticated");
        }

        // Extract and return the UserDto from the authentication
        log.debug("User authentication found, extracting UserDto.");
        return extractUserDto(userAuthentication);
    }

    /**
     * Helper method to extract {@link UserDto} from a {@link UserAuthentication} object.
     * <p>
     * This method maps the {@link User} principal from the {@link UserAuthentication} to a {@link UserDto}.
     * If the {@code principal} is null, it throws an {@link InsufficientAuthenticationException}.
     * </p>
     *
     * @param userAuthentication the {@link UserAuthentication} object from which the user details are extracted.
     *
     * @return the extracted {@link UserDto}.
     *
     * @throws InsufficientAuthenticationException if the {@code principal} is not set.
     */
    private UserDto extractUserDto(UserAuthentication userAuthentication) {
        log.debug("Extracting UserDto from UserAuthentication.");

        // Attempt to map User to UserDto, throwing an exception if the user is not set.
        return Optional.ofNullable(userAuthentication.getPrincipal()).map(u -> {
            log.info("User found: {}. Extracting details into UserDto.", u.getEmail());
            return UserDto.builder()
                          .userId(u.getUserId())
                          .email(u.getEmail())
                          .jwtToken(userAuthentication.getJwtToken())
                          .build();
        }).orElseThrow(() -> {
            log.error("Failed to extract user from authentication: User is not set.");
            return new InsufficientAuthenticationException("User is not set to Authentication");
        });
    }
}
