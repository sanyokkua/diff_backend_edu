import { createAsyncThunk }                                                                    from "@reduxjs/toolkit";
import { parseErrorMessage, ResponseDto, UserClient, UserDeletionDTO, UserDto, UserUpdateDTO } from "../../api";
import { AxiosClient, LogLevel }                                                               from "../../config";


const logger = LogLevel.getLogger("UserThunks");

/**
 * Checks if the response is a valid user response.
 * @function isValidUserResponse
 * @param {ResponseDto<UserDto> | null | undefined} response - The response to check.
 * @returns {response is ResponseDto<UserDto>} - True if the response is valid, false otherwise.
 */
const isValidUserResponse = (response: ResponseDto<UserDto> | null | undefined): response is ResponseDto<UserDto> => {
    return !!(response && response.data && response.data.userId && response.data.email);
};

/**
 * Payload for updating a user.
 * @property {number} userId - The ID of the user.
 * @property {UserUpdateDTO} userUpdateDTO - The data for updating the user.
 */
export type UserUpdateRequest = {
    userId: number;
    userUpdateDTO: UserUpdateDTO;
};

/**
 * Payload for deleting a user.
 * @property {number} userId - The ID of the user.
 * @property {UserDeletionDTO} userDeletionDTO - The data for deleting the user.
 */
export type UserDeleteRequest = {
    userId: number;
    userDeletionDTO: UserDeletionDTO;
};

/**
 * Thunk to fetch a user.
 * @function fetchUser
 * @param {number} userId - The ID of the user to fetch.
 * @param {Object} thunkAPI - The thunk API object.
 * @param {Function} thunkAPI.rejectWithValue - Function to reject with a value.
 * @returns {Promise<UserDto | string>} - The fetched user data or an error message.
 */
export const fetchUser = createAsyncThunk(
    "users/fetchUser",
    async (userId: number, { rejectWithValue }) => {
        logger.debug(`Attempting to fetch user with ID: ${ userId }`);
        const userClient = new UserClient(AxiosClient);

        try {
            const response = await userClient.getUser(userId);
            if (!isValidUserResponse(response)) {
                const validationError = "Invalid user response structure.";
                logger.warn(`Failed to fetch user with ID: ${ userId } - ${ validationError }`);
                return rejectWithValue(validationError);
            }

            logger.info(`Successfully fetched user with ID: ${ userId }`);
            return response.data;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Failed to fetch user.");
            logger.error(`Error fetching user with ID: ${ userId } - ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);

/**
 * Thunk to update a user.
 * @function updateUser
 * @param {UserUpdateRequest} payload - The payload containing userId and userUpdateDTO.
 * @param {Object} thunkAPI - The thunk API object.
 * @param {Function} thunkAPI.rejectWithValue - Function to reject with a value.
 * @returns {Promise<UserDto | string>} - The updated user data or an error message.
 */
export const updateUser = createAsyncThunk(
    "users/updateUser",
    async ({ userId, userUpdateDTO }: UserUpdateRequest, { rejectWithValue }) => {
        logger.debug(`Attempting to update user with ID: ${ userId }`);
        const userClient = new UserClient(AxiosClient);

        try {
            const response = await userClient.updateUser(userId, userUpdateDTO);
            if (!isValidUserResponse(response)) {
                const validationError = "Invalid user response structure.";
                logger.warn(`Failed to update user with ID: ${ userId } - ${ validationError }`);
                return rejectWithValue(validationError);
            }

            logger.info(`Successfully updated user with ID: ${ userId }`);
            return response.data;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Failed to update user.");
            logger.error(`Error updating user with ID: ${ userId } - ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);

/**
 * Thunk to delete a user.
 * @function deleteUser
 * @param {UserDeleteRequest} payload - The payload containing userId and userDeletionDTO.
 * @param {Object} thunkAPI - The thunk API object.
 * @param {Function} thunkAPI.rejectWithValue - Function to reject with a value.
 * @returns {Promise<null | string>} - Null if the user was deleted successfully, or an error message.
 */
export const deleteUser = createAsyncThunk(
    "users/deleteUser",
    async ({ userId, userDeletionDTO }: UserDeleteRequest, { rejectWithValue }) => {
        logger.debug(`Attempting to delete user with ID: ${ userId }`);
        const userClient = new UserClient(AxiosClient);

        try {
            const response = await userClient.deleteUser(userId, userDeletionDTO);
            logger.info(`Successfully deleted user with ID: ${ userId }`);
            return response.data;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Failed to delete user.");
            logger.error(`Error deleting user with ID: ${ userId } - ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);
