package ua.kostenko.tasks.app.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import ua.kostenko.tasks.app.dto.user.UserCreationDTO;
import ua.kostenko.tasks.app.dto.user.UserDto;
import ua.kostenko.tasks.app.dto.user.UserLoginDto;
import ua.kostenko.tasks.app.entity.User;
import ua.kostenko.tasks.app.exception.InvalidPasswordException;
import ua.kostenko.tasks.app.repository.UserRepository;
import ua.kostenko.tasks.app.utility.UserUtils;

/**
 * Service class responsible for handling user authentication and registration.
 * <p>
 * This class provides functionality for logging in users, validating credentials, generating JWT tokens,
 * and registering new users.
 * </p>
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class AuthenticationService {

    /**
     * Service for handling user-related operations like creation and retrieval.
     */
    private final UserService userService;

    /**
     * Repository for accessing user data from the database.
     */
    private final UserRepository userRepository;

    /**
     * Service for generating and validating JWT tokens.
     */
    private final JwtService jwtService;

    /**
     * Password encoder for securely comparing and encoding passwords.
     */
    private final PasswordEncoder passwordEncoder;

    /**
     * Logs in a user by validating their email and password, and generates a JWT token upon successful login.
     *
     * @param userLoginDto the user's login credentials (email and password)
     *
     * @return {@link UserDto} containing the user details and a JWT token
     *
     * @throws IllegalArgumentException if the user with the provided email is not found
     * @throws InvalidPasswordException if the password is incorrect
     */
    public UserDto loginUser(UserLoginDto userLoginDto) {
        UserUtils.validateUserLoginDto(userLoginDto);
        log.info("Attempting to log in user with email: {}", userLoginDto.getEmail());

        // Fetch user by email
        User user = userRepository.findByEmail(userLoginDto.getEmail()).orElseThrow(() -> {
            log.warn("User not found for email: {}", userLoginDto.getEmail());
            return new IllegalArgumentException("User not found");
        });

        // Validate password
        if (!passwordEncoder.matches(userLoginDto.getPassword(), user.getPasswordHash())) {
            log.warn("Invalid password for user with email: {}", userLoginDto.getEmail());
            throw new InvalidPasswordException("Invalid credentials");
        }

        // Generate JWT token
        String jwtToken = jwtService.generateJwtToken(user.getEmail());
        log.info("JWT token generated for user with email: {}", user.getEmail());

        // Return UserDto with token
        return UserDto.builder()
                      .userId(user.getUserId())
                      .email(user.getEmail())
                      .jwtToken(jwtToken)  // Include JWT token in the response
                      .build();
    }

    /**
     * Registers a new user by creating the user in the system and generating a JWT token upon successful registration.
     *
     * @param userCreationDTO the new user's registration details (email, password, etc.)
     *
     * @return {@link UserDto} containing the new user's details and a JWT token
     */
    public UserDto registerUser(UserCreationDTO userCreationDTO) {
        UserUtils.validateUserCreationDTO(userCreationDTO);
        log.info("Registering user with email: {}", userCreationDTO.getEmail());

        // Validate and create user
        UserDto newUser = userService.createUser(userCreationDTO);
        log.info("User successfully created with email: {}", newUser.getEmail());

        // Generate JWT token for the newly registered user
        String jwtToken = jwtService.generateJwtToken(newUser.getEmail());
        log.info("JWT token generated for new user with email: {}", newUser.getEmail());

        // Return UserDto with token
        return UserDto.builder()
                      .userId(newUser.getUserId())
                      .email(newUser.getEmail())
                      .jwtToken(jwtToken)  // Include JWT token in the response
                      .build();
    }
}
