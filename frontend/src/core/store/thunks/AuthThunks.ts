import { createAsyncThunk }                          from "@reduxjs/toolkit";
import { AuthClient, UserCreationDTO, UserLoginDto } from "../../api";
import { parseErrorMessage }                         from "../../api/client/Utils";
import { AxiosClient, LogLevel }                     from "../../config";


const log = LogLevel.getLogger("AuthThunks");

export const loginUser = createAsyncThunk(
    "globals/loginUser",
    async (loginRequest: UserLoginDto, { rejectWithValue }) => {
        log.debug(`Attempting to log in user with email: ${ loginRequest.email }`);
        const client = new AuthClient(AxiosClient);
        try {
            const response = await client.loginUser(loginRequest);
            log.info(`Successfully logged in user with email: ${ loginRequest.email }`);
            return response;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Failed to login");
            log.error(`Error logging in user with email: ${ loginRequest.email } - ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);

export const registerUser = createAsyncThunk(
    "globals/registerUser",
    async (registerRequest: UserCreationDTO, { rejectWithValue }) => {
        log.debug(`Attempting to register user with email: ${ registerRequest.email }`);
        const client = new AuthClient(AxiosClient);
        try {
            const response = await client.registerUser(registerRequest);
            log.info(`Successfully registered user with email: ${ registerRequest.email }`);
            return response;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Failed to Register");
            log.error(`Error registering user with email: ${ registerRequest.email } - ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);
