package ua.kostenko.tasks.app.service;

import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import ua.kostenko.tasks.app.dto.user.UserCreationDTO;
import ua.kostenko.tasks.app.dto.user.UserDeletionDTO;
import ua.kostenko.tasks.app.dto.user.UserDto;
import ua.kostenko.tasks.app.dto.user.UserUpdateDTO;
import ua.kostenko.tasks.app.entity.User;
import ua.kostenko.tasks.app.exception.EmailAlreadyExistsException;
import ua.kostenko.tasks.app.exception.InvalidEmailFormatException;
import ua.kostenko.tasks.app.exception.InvalidPasswordException;
import ua.kostenko.tasks.app.repository.UserRepository;
import ua.kostenko.tasks.app.utility.UserUtils;

import java.util.Optional;

/**
 * Service class for handling user-related operations such as creation, updating, and deletion of users.
 */
@Slf4j
@Service
@RequiredArgsConstructor
@Transactional
public class UserService {

    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;

    /**
     * Creates a new user.
     *
     * @param userCreationDTO the DTO containing user creation data
     *
     * @return the created UserDto
     *
     * @throws InvalidEmailFormatException if the email format is invalid
     * @throws EmailAlreadyExistsException if the email already exists in the repository
     * @throws InvalidPasswordException    if passwords do not match
     */
    public UserDto createUser(UserCreationDTO userCreationDTO) {
        UserUtils.validateUserCreationDTO(userCreationDTO);
        log.info("Creating user with email: {}", userCreationDTO.getEmail());

        UserUtils.validateEmailFormat(userCreationDTO.getEmail());
        UserUtils.checkUserExists(userRepository, userCreationDTO.getEmail());
        UserUtils.validatePasswords(userCreationDTO.getPassword(), userCreationDTO.getPasswordConfirmation());

        User newUser = User.builder()
                           .email(userCreationDTO.getEmail())
                           .passwordHash(passwordEncoder.encode(userCreationDTO.getPassword()))
                           .build();

        User savedUser = userRepository.save(newUser);
        log.info("User created successfully with email: {}", savedUser.getEmail());

        return UserDto.builder().userId(savedUser.getUserId()).email(savedUser.getEmail()).build();
    }

    /**
     * Updates the user's password.
     *
     * @param userId        the ID of the user
     * @param userUpdateDTO the DTO containing current and new password information
     *
     * @return the updated UserDto
     *
     * @throws InvalidPasswordException if the current password is incorrect or the password update validation fails
     */
    public UserDto updateUserPassword(Long userId, UserUpdateDTO userUpdateDTO) {
        UserUtils.validateUserUpdateDTO(userUpdateDTO);
        log.info("Updating password for user with ID: {}", userId);

        User user = userRepository.findById(userId).orElseThrow(() -> new IllegalArgumentException("User not found"));

        UserUtils.validatePasswordUpdate(userUpdateDTO, user, passwordEncoder);

        user.setPasswordHash(passwordEncoder.encode(userUpdateDTO.getNewPassword()));
        User updatedUser = userRepository.save(user);

        log.info("Password updated successfully for user with ID: {}", updatedUser.getUserId());
        return UserDto.builder().userId(updatedUser.getUserId()).email(updatedUser.getEmail()).build();
    }

    /**
     * Deletes a user.
     *
     * @param userDeletionDTO the DTO containing user ID and current password
     *
     * @throws InvalidPasswordException if the current password is incorrect
     * @throws IllegalArgumentException if the user is not found
     */
    public void deleteUser(Long userId, UserDeletionDTO userDeletionDTO) {
        UserUtils.validateUserDeletionDTO(userDeletionDTO);
        log.info("Deleting user with ID: {}", userId);

        Optional<User> userOptional = userRepository.findById(userId);

        if (userOptional.isEmpty()) {
            log.warn("User not found with ID: {}", userId);
            throw new IllegalArgumentException("User not found");
        }

        User user = userOptional.get();

        if (!passwordEncoder.matches(userDeletionDTO.getCurrentPassword(), user.getPasswordHash())) {
            log.warn("Invalid current password for user with ID: {}", user.getUserId());
            throw new InvalidPasswordException("Invalid current password");
        }

        userRepository.delete(user);
        log.info("User deleted successfully with ID: {}", user.getUserId());
    }
}
