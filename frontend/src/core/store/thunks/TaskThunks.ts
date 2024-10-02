import { createAsyncThunk }                           from "@reduxjs/toolkit";
import { TaskClient, TaskCreationDTO, TaskUpdateDTO } from "../../api";
import { parseErrorMessage }                          from "../../api/client/Utils";
import { AxiosClient, LogLevel }                      from "../../config";


const log = LogLevel.getLogger("TaskThunks");

export type GetUserTaskRequest = {
    userId: number;
    taskId: number;
};
export type CreateTaskRequest = {
    userId: number;
    taskCreationDTO: TaskCreationDTO;
};
export type UpdateTaskRequest = {
    userId: number;
    taskId: number;
    taskUpdateDTO: TaskUpdateDTO;
};
export type DeleteUserTaskRequest = GetUserTaskRequest;

export const getTask = createAsyncThunk(
    "globals/getTask",
    async ({ userId, taskId }: GetUserTaskRequest, { rejectWithValue }) => {
        log.debug(`Fetching task with ID: ${ taskId } for user ID: ${ userId }`);
        const client = new TaskClient(AxiosClient);
        try {
            const response = await client.getTask(userId, taskId);
            log.info(`Successfully fetched task with ID: ${ taskId } for user ID: ${ userId }`);
            return response;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Failed to fetch tasks");
            log.error(`Error fetching task with ID: ${ taskId } for user ID: ${ userId } - ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);

export const getTasks = createAsyncThunk(
    "globals/getTasks",
    async (userId: number, { rejectWithValue }) => {
        log.debug(`Fetching tasks for user ID: ${ userId }`);
        const client = new TaskClient(AxiosClient);
        try {
            const response = await client.getTasks(userId);
            log.info(`Successfully fetched tasks for user ID: ${ userId }`);
            return response;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Failed to fetch tasks");
            log.error(`Error fetching tasks for user ID: ${ userId } - ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);

export const createTask = createAsyncThunk(
    "globals/createTask",
    async ({ userId, taskCreationDTO }: CreateTaskRequest, { rejectWithValue }) => {
        log.debug(`Creating task for user ID: ${ userId }`);
        const client = new TaskClient(AxiosClient);
        try {
            const response = await client.createTask(userId, taskCreationDTO);
            log.info(`Successfully created task for user ID: ${ userId }`);
            return response;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Failed to create task");
            log.error(`Error creating task for user ID: ${ userId } - ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);

export const updateTask = createAsyncThunk(
    "globals/updateTask",
    async ({ userId, taskId, taskUpdateDTO }: UpdateTaskRequest, { rejectWithValue }) => {
        log.debug(`Updating task with ID: ${ taskId } for user ID: ${ userId }`);
        const client = new TaskClient(AxiosClient);
        try {
            const response = await client.updateTask(userId, taskId, taskUpdateDTO);
            log.info(`Successfully updated task with ID: ${ taskId } for user ID: ${ userId }`);
            return response;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Failed to update task");
            log.error(`Error updating task with ID: ${ taskId } for user ID: ${ userId } - ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);

export const deleteTask = createAsyncThunk(
    "globals/deleteTask",
    async ({ userId, taskId }: DeleteUserTaskRequest, { rejectWithValue }) => {
        log.debug(`Deleting task with ID: ${ taskId } for user ID: ${ userId }`);
        const client = new TaskClient(AxiosClient);
        try {
            const response = await client.deleteTask(userId, taskId);
            log.info(`Successfully deleted task with ID: ${ taskId } for user ID: ${ userId }`);
            return response;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Failed to delete task");
            log.error(`Error deleting task with ID: ${ taskId } for user ID: ${ userId } - ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);
