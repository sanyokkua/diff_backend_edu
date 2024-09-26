package ua.kostenko.tasks.app.controller;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import ua.kostenko.tasks.app.dto.ResponseDto;
import ua.kostenko.tasks.app.dto.user.UserCreationDTO;
import ua.kostenko.tasks.app.dto.user.UserDto;
import ua.kostenko.tasks.app.dto.user.UserLoginDto;
import ua.kostenko.tasks.app.service.AuthenticationService;
import ua.kostenko.tasks.app.utility.ResponseDtoUtils;

/**
 * REST controller for handling authentication-related requests.
 * <p>
 * This controller provides endpoints for user login and registration.
 * </p>
 */
@Slf4j
@RestController
@RequestMapping("api/v1/auth")
@RequiredArgsConstructor
@Tag(name = "Authentication REST Controller", description = "Handles user login and registration requests.")
public class AuthController {

    private final AuthenticationService authenticationService;

    /**
     * Handles user login requests.
     *
     * @param userLoginDto the user login data transfer object containing email and password
     *
     * @return a ResponseEntity containing the ResponseDto with user information or an error message
     */
    @PostMapping("/login")
    @Operation(summary = "User Login", description = "Handles user login by validating the provided email and password.")
    public ResponseEntity<ResponseDto<UserDto>> loginUser(
            @RequestBody @Parameter(description = "User login details including email and password.") UserLoginDto userLoginDto) {
        log.info("Received login request for email: {}", userLoginDto.getEmail());

        UserDto userDto = authenticationService.loginUser(userLoginDto);
        log.info("Login successful for email: {}", userLoginDto.getEmail());
        return ResponseDtoUtils.buildDtoResponse(userDto, HttpStatus.OK);
    }

    /**
     * Handles user registration requests.
     *
     * @param userCreationDTO the user creation data transfer object containing user details
     *
     * @return a ResponseEntity containing the ResponseDto with user information or an error message
     */
    @PostMapping("/register")
    @Operation(summary = "User Registration", description = "Handles user registration by creating a new user account with the provided details.")
    public ResponseEntity<ResponseDto<UserDto>> registerUser(
            @RequestBody @Parameter(description = "User creation details including email, password, and confirmation.") UserCreationDTO userCreationDTO) {
        log.info("Received registration request for email: {}", userCreationDTO.getEmail());

        UserDto userDto = authenticationService.registerUser(userCreationDTO);
        log.info("Registration successful for email: {}", userCreationDTO.getEmail());
        return ResponseDtoUtils.buildDtoResponse(userDto, HttpStatus.CREATED);
    }
}
