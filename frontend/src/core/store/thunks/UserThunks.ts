import { createAsyncThunk }                           from "@reduxjs/toolkit";
import { UserClient, UserDeletionDTO, UserUpdateDTO } from "../../api";
import { parseErrorMessage }                          from "../../api/client/Utils";
import { AxiosClient, LogLevel }                      from "../../config";


const log = LogLevel.getLogger("UserThunks");

export type UserUpdateRequest = {
    userId: number;
    userUpdateDTO: UserUpdateDTO;
};

export type UserDeleteRequest = {
    userId: number;
    userDeletionDTO: UserDeletionDTO;
};

export const getUser = createAsyncThunk(
    "globals/getUser",
    async (userId: number, { rejectWithValue }) => {
        log.debug(`Fetching user with ID: ${ userId }`);
        const client = new UserClient(AxiosClient);
        try {
            const response = await client.getUser(userId);
            log.info(`Successfully fetched user with ID: ${ userId }`);
            return response;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Failed to fetch user");
            log.error(`Error fetching user with ID: ${ userId } - ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);

export const updateUser = createAsyncThunk(
    "globals/updateUser",
    async ({ userId, userUpdateDTO }: UserUpdateRequest, { rejectWithValue }) => {
        log.debug(`Updating user with ID: ${ userId }`);
        const client = new UserClient(AxiosClient);
        try {
            const response = await client.updateUser(userId, userUpdateDTO);
            log.info(`Successfully updated user with ID: ${ userId }`);
            return response;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Failed to update user");
            log.error(`Error updating user with ID: ${ userId } - ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);

export const deleteUser = createAsyncThunk(
    "globals/deleteUser",
    async ({ userId, userDeletionDTO }: UserDeleteRequest, { rejectWithValue }) => {
        log.debug(`Deleting user with ID: ${ userId }`);
        const client = new UserClient(AxiosClient);
        try {
            const response = await client.deleteUser(userId, userDeletionDTO);
            log.info(`Successfully deleted user with ID: ${ userId }`);
            return response;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Failed to delete user");
            log.error(`Error deleting user with ID: ${ userId } - ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);
