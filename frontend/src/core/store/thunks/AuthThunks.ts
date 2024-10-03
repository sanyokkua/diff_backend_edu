import { createAsyncThunk }                                                                   from "@reduxjs/toolkit";
import { AuthClient, parseErrorMessage, ResponseDto, UserCreationDTO, UserDto, UserLoginDto } from "../../api";
import { AxiosClient, LogLevel }                                                              from "../../config";


const logger = LogLevel.getLogger("AuthThunks");

/**
 * Checks if the response is a valid authentication response.
 * @function isValidAuthResponse
 * @param {ResponseDto<UserDto> | undefined | null} response - The response to check.
 * @returns {response is ResponseDto<UserDto>} - True if the response is valid, false otherwise.
 */
const isValidAuthResponse = (response: ResponseDto<UserDto> | undefined | null): response is ResponseDto<UserDto> => {
    return !!(response && response.data && response.data.userId && response.data.email && response.data.jwtToken);
};

/**
 * Thunk for logging in a user.
 * @function loginUser
 * @param {UserLoginDto} loginData - The login data.
 * @param {Object} thunkAPI - The thunk API object.
 * @param {Function} thunkAPI.rejectWithValue - Function to reject with a value.
 * @returns {Promise<UserDto | string>} - The logged-in user data or an error message.
 */
export const loginUser = createAsyncThunk(
    "auth/loginUser",
    async (loginData: UserLoginDto, { rejectWithValue }) => {
        logger.debug(`Login attempt for email: ${ loginData.email }`);
        const authClient = new AuthClient(AxiosClient);

        try {
            // Let TypeScript infer the type of the response
            const response = await authClient.loginUser(loginData);

            if (!isValidAuthResponse(response)) {
                const validationError = "Invalid response structure from login API.";
                logger.warn(`Login failed for email: ${ loginData.email } - ${ validationError }`);
                return rejectWithValue(validationError);
            }

            logger.info(`Login successful for user: ${ response.data.email }`);
            return response.data;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Login failed.");
            logger.error(`Login error for email: ${ loginData.email } - ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);

/**
 * Thunk for registering a user.
 * @function registerUser
 * @param {UserCreationDTO} registrationData - The registration data.
 * @param {Object} thunkAPI - The thunk API object.
 * @param {Function} thunkAPI.rejectWithValue - Function to reject with a value.
 * @returns {Promise<UserDto | string>} - The registered user data or an error message.
 */
export const registerUser = createAsyncThunk(
    "auth/registerUser",
    async (registrationData: UserCreationDTO, { rejectWithValue }) => {
        logger.debug(`Registration attempt for email: ${ registrationData.email }`);
        const authClient = new AuthClient(AxiosClient);

        try {
            const response = await authClient.registerUser(registrationData);

            if (!isValidAuthResponse(response)) {
                const validationError = "Invalid response structure from registration API.";
                logger.warn(`Registration failed for email: ${ registrationData.email } - ${ validationError }`);
                return rejectWithValue(validationError);
            }

            logger.info(`Registration successful for user: ${ response.data.email }`);
            return response.data;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Registration failed.");
            logger.error(`Registration error for email: ${ registrationData.email } - ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);
