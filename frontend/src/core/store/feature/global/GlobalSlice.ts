import { createSlice, PayloadAction }                        from "@reduxjs/toolkit";
import { jwtDecode, JwtPayload }                             from "jwt-decode";
import { getDateFromSeconds, ResponseDto, TaskDto, UserDto } from "../../../api";

import { LogLevel }               from "../../../config";
import BrowserStore, { ItemType } from "../../BrowserStore";

import {
    createTask,
    deleteTask,
    deleteUser,
    getTask,
    getTasks,
    getUser,
    loginUser,
    registerUser,
    updateTask,
    updateUser
} from "../../thunks";


const log = LogLevel.getLogger("GlobalSlice");

export interface GlobalState {
    userJwtToken: string;
    userEmail: string;
    userId: number;
    userIsLoggedIn: boolean;
    taskId: number;
    taskName: string;
    taskDescription: string;
    tasksList: TaskDto[];
    appBarHeader: string;
    appIsLoading: boolean;
    appError: string;
}

const initializeState = (): GlobalState => {
    const storeApi = new BrowserStore();

    const userId = parseInt(storeApi.getData(ItemType.USER_ID) ?? "0");
    const userEmail = storeApi.getData(ItemType.USER_EMAIL) ?? "";
    const userJwt = storeApi.getData(ItemType.JWT_TOKEN) ?? "";

    let userIsLoggedIn = false;
    let error = "";

    if (userJwt && userId && userEmail) {
        try {
            const decoded = jwtDecode<JwtPayload>(userJwt);
            if (decoded?.sub && decoded?.exp) {
                const userTimeExp = decoded.exp;
                userIsLoggedIn = new Date().getTime() < getDateFromSeconds(userTimeExp).getTime();
            } else {
                error = "Failed to decode token";
            }
        } catch (e) {
            log.warn("Failed to parse JWT token.", e);
            error = "Failed to load token from the history";
        }
    }

    return {
        userJwtToken: userJwt,
        userEmail: userEmail,
        userId: userId,
        userIsLoggedIn: userIsLoggedIn,
        taskId: 0,
        taskName: "",
        taskDescription: "",
        tasksList: [],
        appBarHeader: "",
        appIsLoading: false,
        appError: error
    };
};

const InitialAppState: GlobalState = initializeState();

const resetUserDetails = (state: GlobalState): void => {
    state.userId = 0;
    state.userEmail = "";
    state.userJwtToken = "";
    state.userIsLoggedIn = false;
};

const loginOrRegisterFulFilled = (action: PayloadAction<ResponseDto<UserDto>>, state: GlobalState): void => {
    const payload = action.payload?.data;
    if (!payload?.jwtToken) {
        resetUserDetails(state);
        state.appError = "No JWT token in response";
        return;
    }

    state.userJwtToken = payload.jwtToken;
    state.userEmail = payload.email;
    state.userId = payload.userId;
    state.userIsLoggedIn = true;
    state.appError = "";
};

export const globalSlice = createSlice(
    {
        name: "globals",
        initialState: InitialAppState,
        reducers: {
            setUserJwtToken: (state, action: PayloadAction<string>) => {
                log.debug(`Setting current user JWT: ${ action.payload }`);
                state.userJwtToken = action.payload;
            },
            setUserEmail: (state, action: PayloadAction<string>) => {
                log.debug(`Setting current user Email: ${ action.payload }`);
                state.userEmail = action.payload;
            },
            setUserId: (state, action: PayloadAction<number>) => {
                log.debug(`Setting current user ID: ${ action.payload }`);
                state.userId = action.payload;
            },
            setUserIsLoggedIn: (state, action: PayloadAction<boolean>) => {
                log.debug(`Setting current user isLoggged: ${ action.payload }`);
                state.userIsLoggedIn = action.payload;
            },
            setTaskId: (state, action: PayloadAction<number>) => {
                log.debug(`Setting current task ID: ${ action.payload }`);
                state.taskId = action.payload;
            },
            setTaskName: (state, action: PayloadAction<string>) => {
                log.debug(`Setting current task name: ${ action.payload }`);
                state.taskName = action.payload;
            },
            setTaskDescription: (state, action: PayloadAction<string>) => {
                log.debug(`Setting current task description: ${ action.payload }`);
                state.taskDescription = action.payload;
            },
            setTasksList: (state, action: PayloadAction<TaskDto[]>) => {
                log.debug(`Setting tasks: ${ action.payload }`);
                state.tasksList = action.payload;
            },
            setAppBarHeader: (state, action: PayloadAction<string>) => {
                log.debug(`Setting current app bar header: ${ action.payload }`);
                state.appBarHeader = action.payload;
            },
            setAppIsLoading: (state, action: PayloadAction<boolean>) => {
                log.debug(`Setting current app is loading: ${ action.payload }`);
                state.appIsLoading = action.payload;
            },
            setAppError: (state, action: PayloadAction<string>) => {
                log.debug(`Setting current error: ${ action.payload }`);
                state.appError = action.payload;
            }
        },
        extraReducers: (builder) => {
            builder
                // User-related actions
                .addCase(loginUser.pending, (state: GlobalState) => {
                    log.debug("Login request pending");
                    resetUserDetails(state);
                    state.appIsLoading = true;
                })
                .addCase(loginUser.fulfilled, (state: GlobalState, action: PayloadAction<ResponseDto<UserDto>>) => {
                    log.info("Login request fulfilled successfully");
                    loginOrRegisterFulFilled(action, state);
                    state.appIsLoading = false;
                })
                .addCase(loginUser.rejected, (state: GlobalState, action) => {
                    resetUserDetails(state);
                    state.appIsLoading = false;
                    state.appError = action.payload as string || "Login request was rejected";
                })

                .addCase(registerUser.pending, (state: GlobalState) => {
                    log.debug("Register request pending");
                    resetUserDetails(state);
                    state.appIsLoading = true;
                })
                .addCase(registerUser.fulfilled, (state: GlobalState, action: PayloadAction<ResponseDto<UserDto>>) => {
                    log.info("Register request fulfilled successfully");
                    loginOrRegisterFulFilled(action, state);
                })
                .addCase(registerUser.rejected, (state: GlobalState, action) => {
                    resetUserDetails(state);
                    state.appIsLoading = false;
                    state.appError = action.payload as string || "Register request was rejected";
                })

                .addCase(getUser.pending, (state: GlobalState) => {
                    log.debug("GetUser request pending");
                    state.userId = 0;
                    state.userEmail = "";
                    state.appIsLoading = true;
                })
                .addCase(getUser.fulfilled, (state: GlobalState, action: PayloadAction<ResponseDto<UserDto>>) => {
                    log.info("GetUser request fulfilled successfully");
                    const payload = action.payload?.data;
                    state.userEmail = payload.email;
                    state.userId = payload.userId;
                    state.appError = "";
                })
                .addCase(getUser.rejected, (state: GlobalState, action) => {
                    state.userId = 0;
                    state.userEmail = "";
                    state.appIsLoading = false;
                    state.appError = action.payload as string || "getUser request was rejected";
                })

                .addCase(updateUser.pending, (state: GlobalState) => {
                    log.debug("updateUser request pending"); // As updating password, everything should be reset
                    state.userId = 0;
                    state.userEmail = "";
                    state.userJwtToken = "";
                    state.userIsLoggedIn = false;
                    state.appIsLoading = true;
                })
                .addCase(updateUser.fulfilled, (state: GlobalState, action: PayloadAction<ResponseDto<UserDto>>) => {
                    log.info("updateUser request fulfilled successfully");
                    const payload = action.payload?.data;
                    if (!payload?.jwtToken) {
                        state.userId = 0;
                        state.userEmail = "";
                        state.userJwtToken = "";
                        state.userIsLoggedIn = false;
                        state.appIsLoading = false;
                        state.appError = "No JWT token in response";
                        return;
                    }

                    state.userJwtToken = payload.jwtToken;
                    state.userEmail = payload.email;
                    state.userId = payload.userId;
                    state.userIsLoggedIn = true;
                    state.appIsLoading = false;
                    state.appError = "";
                })
                .addCase(updateUser.rejected, (state: GlobalState, action) => {
                    state.userId = 0;
                    state.userEmail = "";
                    state.userJwtToken = "";
                    state.userIsLoggedIn = false;
                    state.appIsLoading = false;
                    state.appError = action.payload as string || "updateUser request was rejected";
                })

                .addCase(deleteUser.pending, (state: GlobalState) => {
                    log.debug("deleteUser request pending");
                    state.userId = 0;
                    state.userEmail = "";
                    state.userJwtToken = "";
                    state.userIsLoggedIn = false;
                    state.appIsLoading = true;
                })
                .addCase(deleteUser.fulfilled, (state: GlobalState) => {
                    log.info("deleteUser request fulfilled successfully");
                    state.userId = 0;
                    state.userEmail = "";
                    state.userJwtToken = "";
                    state.userIsLoggedIn = false;
                    state.appIsLoading = false;
                    state.appError = "";
                })
                .addCase(deleteUser.rejected, (state: GlobalState, action) => {
                    state.userId = 0;
                    state.userEmail = "";
                    state.userJwtToken = "";
                    state.userIsLoggedIn = false;
                    state.appIsLoading = false;
                    state.appError = action.payload as string || "deleteUser request was rejected";
                })

                // Task-related actions
                .addCase(getTask.pending, (state: GlobalState) => {
                    log.debug("getTask request pending");
                    state.appIsLoading = true;
                    state.taskId = 0;
                    state.taskName = "";
                    state.taskDescription = "";
                })
                .addCase(getTask.fulfilled, (state: GlobalState, action: PayloadAction<ResponseDto<TaskDto>>) => {
                    log.info("getTask request fulfilled successfully");
                    const payload = action.payload?.data;
                    state.taskId = payload?.taskId ?? 0;
                    state.taskName = payload?.name ?? "";
                    state.taskDescription = payload?.description ?? "";

                    state.appIsLoading = false;
                    state.appError = "";
                })
                .addCase(getTask.rejected, (state: GlobalState, action) => {
                    state.taskId = 0;
                    state.taskName = "";
                    state.taskDescription = "";
                    state.appIsLoading = false;
                    state.appError = action.payload as string || "getTask request was rejected";
                })

                .addCase(getTasks.pending, (state: GlobalState) => {
                    log.debug("getTasks request pending");
                    state.tasksList = [];
                    state.appIsLoading = true;
                })
                .addCase(getTasks.fulfilled, (state: GlobalState, action: PayloadAction<ResponseDto<TaskDto[]>>) => {
                    log.info("getTasks request fulfilled successfully");
                    const payload = action.payload?.data;
                    state.tasksList = payload ?? [];
                    state.appIsLoading = false;
                })
                .addCase(getTasks.rejected, (state: GlobalState, action) => {
                    state.tasksList = [];
                    state.appIsLoading = false;
                    state.appError = action.payload as string || "getTasks request was rejected";
                })

                .addCase(createTask.pending, (state: GlobalState) => {
                    log.debug("createTask request pending");
                    state.taskId = 0;
                    state.taskName = "";
                    state.appIsLoading = true;
                })
                .addCase(createTask.fulfilled, (state: GlobalState, action: PayloadAction<ResponseDto<TaskDto>>) => {
                    log.info("createTask request fulfilled successfully");
                    const payload = action.payload?.data;
                    state.taskId = payload?.taskId ?? 0;
                    state.taskName = payload?.name ?? "";
                    state.appIsLoading = false;
                    state.appError = "";
                })
                .addCase(createTask.rejected, (state: GlobalState, action) => {
                    state.taskId = 0;
                    state.taskName = "";
                    state.appIsLoading = false;
                    state.appError = action.payload as string || "createTask request was rejected";
                })

                .addCase(updateTask.pending, (state: GlobalState) => {
                    log.debug("updateTask request pending");
                    state.taskId = 0;
                    state.taskName = "";
                    state.appIsLoading = true;
                })
                .addCase(updateTask.fulfilled, (state: GlobalState, action: PayloadAction<ResponseDto<TaskDto>>) => {
                    log.info("updateTask request fulfilled successfully");
                    const payload = action.payload?.data;
                    state.taskId = payload?.taskId ?? 0;
                    state.taskName = payload?.name ?? "";
                    state.appIsLoading = false;
                    state.appError = "";
                })
                .addCase(updateTask.rejected, (state: GlobalState, action) => {
                    state.taskId = 0;
                    state.taskName = "";
                    state.appIsLoading = false;
                    state.appError = action.payload as string || "updateTask request was rejected";
                })

                .addCase(deleteTask.pending, (state: GlobalState) => {
                    log.debug("deleteTask request pending");
                    state.taskId = 0;
                    state.taskName = "";
                    state.taskDescription = "";
                    state.appIsLoading = true;
                })
                .addCase(deleteTask.fulfilled, (state: GlobalState) => {
                    log.info("deleteTask request fulfilled successfully");
                    state.taskId = 0;
                    state.taskName = "";
                    state.taskDescription = "";
                    state.appIsLoading = false;
                    state.appError = "";
                })
                .addCase(deleteTask.rejected, (state: GlobalState, action) => {
                    state.taskId = 0;
                    state.taskName = "";
                    state.taskDescription = "";
                    state.appIsLoading = false;
                    state.appError = action.payload as string || "deleteTask request was rejected";
                })
            ;
        }
    }
);

export const {
    setUserJwtToken,
    setUserEmail,
    setUserId,
    setUserIsLoggedIn,
    setTaskId,
    setTaskName,
    setTaskDescription,
    setTasksList,
    setAppBarHeader,
    setAppIsLoading,
    setAppError
} = globalSlice.actions;
export default globalSlice.reducer;