import { createSlice, PayloadAction }                                 from "@reduxjs/toolkit";
import { jwtDecode, JwtPayload }                                      from "jwt-decode";
import { getDateFromSeconds, parseErrorMessage, UserDto }             from "../../../api";
import { LogLevel }                                                   from "../../../config";
import BrowserStore, { ItemType }                                     from "../../BrowserStore";
import { deleteUser, fetchUser, loginUser, registerUser, updateUser } from "../../thunks";


const logger = LogLevel.getLogger("UserSlice");

/**
 * Represents the state of the user in the application.
 * @interface UserState
 * @property {number} userId - The unique identifier for the user.
 * @property {string} userEmail - The email address of the user.
 * @property {string} userJwtToken - The JWT token for the user session.
 * @property {boolean} userIsLoggedIn - Indicates if the user is logged in.
 * @property {boolean} userIsLoading - Indicates if a user-related action is in progress.
 * @property {string} userError - Error message related to user actions, if any.
 */
export interface UserState {
    userId: number;
    userEmail: string;
    userJwtToken: string;
    userIsLoggedIn: boolean;
    userIsLoading: boolean;
    userError: string;
}

/**
 * Initializes the user state based on stored data and JWT token validation.
 * @returns {UserState} The initial state for the user, populated with stored values and JWT token status.
 */
const initializeState = (): UserState => {
    const storeApi = new BrowserStore();
    const userId = parseInt(storeApi.getData(ItemType.USER_ID) ?? "0");
    const userEmail = storeApi.getData(ItemType.USER_EMAIL) ?? "";
    const userJwtToken = storeApi.getData(ItemType.JWT_TOKEN) ?? "";

    let userIsLoggedIn = false;
    let error = "";

    // Check if the JWT token is valid and extract user information from it
    if (userJwtToken && userId && userEmail) {
        try {
            const decoded = jwtDecode<JwtPayload>(userJwtToken);
            if (decoded?.exp && decoded?.sub) {
                userIsLoggedIn = new Date().getTime() < getDateFromSeconds(decoded.exp).getTime();
            } else {
                error = "JWT token is invalid.";
            }
        } catch (e) {
            logger.warn("Failed to decode JWT token.", e);
            error = parseErrorMessage(e, "Error while loading JWT token.");
        }
    }

    return {
        userId,
        userEmail,
        userJwtToken,
        userIsLoggedIn,
        userIsLoading: false,
        userError: error
    };
};

const initialState: UserState = initializeState();

/**
 * Resets the user details in the state to default values.
 * @param {UserState} state - The current user state to reset.
 * @returns {void}
 */
const resetUserDetails = (state: UserState): void => {
    state.userId = 0;
    state.userEmail = "";
    state.userJwtToken = "";
    state.userIsLoggedIn = false;
};

/**
 * Handles the success of login or registration by updating the state with user information.
 * @param {UserState} state - The current user state to update.
 * @param {UserDto} payload - The data transfer object containing user information.
 * @returns {void}
 */
const handleLoginOrRegisterSuccess = (state: UserState, payload: UserDto): void => {
    if (!payload?.jwtToken) {
        resetUserDetails(state);
        state.userError = "No JWT token in response.";
        return;
    }

    state.userJwtToken = payload.jwtToken;
    state.userEmail = payload.email;
    state.userId = payload.userId;
    state.userIsLoggedIn = true;
    state.userError = "";
};

export const userSlice = createSlice(
    {
        name: "user",
        initialState,
        reducers: {
            /**
             * Sets the JWT token in the user state.
             * @param {UserState} state - The current user state.
             * @param {PayloadAction<string>} action - The action containing the JWT token.
             */
            setUserJwtToken: (state: UserState, action: PayloadAction<string>) => {
                logger.debug(`Setting JWT token: ${ action.payload }`);
                state.userJwtToken = action.payload;
            },
            /**
             * Sets the user's email in the user state.
             * @param {UserState} state - The current user state.
             * @param {PayloadAction<string>} action - The action containing the user's email.
             */
            setUserEmail: (state: UserState, action: PayloadAction<string>) => {
                logger.debug(`Setting user email: ${ action.payload }`);
                state.userEmail = action.payload;
            },
            /**
             * Sets the user's ID in the user state.
             * @param {UserState} state - The current user state.
             * @param {PayloadAction<number>} action - The action containing the user's ID.
             */
            setUserId: (state: UserState, action: PayloadAction<number>) => {
                logger.debug(`Setting user ID: ${ action.payload }`);
                state.userId = action.payload;
            },
            /**
             * Sets the user's login status in the user state.
             * @param {UserState} state - The current user state.
             * @param {PayloadAction<boolean>} action - The action containing the login status.
             */
            setUserIsLoggedIn: (state: UserState, action: PayloadAction<boolean>) => {
                logger.debug(`Setting user login status: ${ action.payload }`);
                state.userIsLoggedIn = action.payload;
            },
            /**
             * Sets the loading state for user actions.
             * @param {UserState} state - The current user state.
             * @param {PayloadAction<boolean>} action - The action containing the loading status.
             */
            setUserIsLoading: (state: UserState, action: PayloadAction<boolean>) => {
                logger.debug(`Setting loading state: ${ action.payload }`);
                state.userIsLoading = action.payload;
            },
            /**
             * Sets an error message in the user state.
             * @param {UserState} state - The current user state.
             * @param {PayloadAction<string>} action - The action containing the error message.
             */
            setUserError: (state: UserState, action: PayloadAction<string>) => {
                logger.debug(`Setting error message: ${ action.payload }`);
                state.userError = action.payload;
            }
        },
        extraReducers: (builder) => {
            builder
                .addCase(loginUser.pending, (state: UserState) => {
                    logger.debug("Login request initiated.");
                    resetUserDetails(state);
                    state.userIsLoading = true;
                })
                .addCase(loginUser.fulfilled, (state: UserState, action: PayloadAction<UserDto>) => {
                    logger.info("Login request succeeded.");
                    handleLoginOrRegisterSuccess(state, action.payload);
                    state.userIsLoading = false;
                })
                .addCase(loginUser.rejected, (state: UserState, action) => {
                    resetUserDetails(state);
                    state.userIsLoading = false;
                    state.userError = action.payload as string || "Login failed.";
                })
                .addCase(registerUser.pending, (state: UserState) => {
                    logger.debug("Registration request initiated.");
                    resetUserDetails(state);
                    state.userIsLoading = true;
                })
                .addCase(registerUser.fulfilled, (state: UserState, action: PayloadAction<UserDto>) => {
                    logger.info("Registration request succeeded.");
                    handleLoginOrRegisterSuccess(state, action.payload);
                    state.userIsLoading = false;
                })
                .addCase(registerUser.rejected, (state: UserState, action) => {
                    resetUserDetails(state);
                    state.userIsLoading = false;
                    state.userError = action.payload as string || "Registration failed.";
                })
                .addCase(fetchUser.pending, (state: UserState) => {
                    logger.debug("Fetching user details.");
                    state.userIsLoading = true;
                })
                .addCase(fetchUser.fulfilled, (state: UserState, action: PayloadAction<UserDto>) => {
                    logger.info("User details fetched successfully.");
                    const { email, userId } = action.payload;
                    state.userEmail = email;
                    state.userId = userId;
                    state.userIsLoading = false;
                    state.userError = "";
                })
                .addCase(fetchUser.rejected, (state: UserState, action) => {
                    logger.error("Fetching user details failed.");
                    state.userIsLoading = false;
                    state.userError = action.payload as string || "Fetch user failed.";
                })
                .addCase(updateUser.pending, (state: UserState) => {
                    logger.debug("Updating user details.");
                    state.userIsLoading = true;
                })
                .addCase(updateUser.fulfilled, (state: UserState, action: PayloadAction<UserDto>) => {
                    logger.info("User details updated successfully.");
                    handleLoginOrRegisterSuccess(state, action.payload);
                    state.userIsLoading = false;
                })
                .addCase(updateUser.rejected, (state: UserState, action) => {
                    logger.error("Updating user details failed.");
                    state.userIsLoading = false;
                    state.userError = action.payload as string || "Update user failed.";
                })
                .addCase(deleteUser.pending, (state: UserState) => {
                    logger.debug("Deleting user.");
                    state.userIsLoading = true;
                })
                .addCase(deleteUser.fulfilled, (state: UserState) => {
                    logger.info("User deleted successfully.");
                    resetUserDetails(state);
                    state.userIsLoading = false;
                    state.userError = "";
                })
                .addCase(deleteUser.rejected, (state: UserState, action) => {
                    logger.error("Deleting user failed.");
                    state.userIsLoading = false;
                    state.userError = action.payload as string || "Delete user failed.";
                });
        }
    }
);

/**
 * User slice action creators.
 *
 * @property {function} setUserJwtToken - Action to set the JWT token.
 * @property {function} setUserEmail - Action to set the user's email.
 * @property {function} setUserId - Action to set the user's ID.
 * @property {function} setUserIsLoggedIn - Action to set the user's login status.
 * @property {function} setUserIsLoading - Action to set the loading status.
 * @property {function} setUserError - Action to set an error message.
 */
export const {
    setUserJwtToken,
    setUserEmail,
    setUserId,
    setUserIsLoggedIn,
    setUserIsLoading,
    setUserError
} = userSlice.actions;

export default userSlice.reducer;
