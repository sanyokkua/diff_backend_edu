import { createAsyncThunk }                                                                    from "@reduxjs/toolkit";
import { parseErrorMessage, ResponseDto, TaskClient, TaskCreationDTO, TaskDto, TaskUpdateDTO } from "../../api";
import { AxiosClient, LogLevel }                                                               from "../../config";


const logger = LogLevel.getLogger("TaskThunks");

/**
 * Payload for fetching a user task.
 * @property {number} userId - The ID of the user.
 * @property {number} taskId - The ID of the task.
 */
export interface GetUserTaskPayload {
    userId: number;
    taskId: number;
}

/**
 * Payload for creating a task.
 * @property {number} userId - The ID of the user.
 * @property {TaskCreationDTO} taskData - The data for the new task.
 */
export interface CreateTaskPayload {
    userId: number;
    taskData: TaskCreationDTO;
}

/**
 * Payload for updating a task.
 * @property {number} userId - The ID of the user.
 * @property {number} taskId - The ID of the task.
 * @property {TaskUpdateDTO} taskData - The updated data for the task.
 */
export interface UpdateTaskPayload {
    userId: number;
    taskId: number;
    taskData: TaskUpdateDTO;
}

/**
 * Payload for deleting a task.
 * @property {number} userId - The ID of the user.
 * @property {number} taskId - The ID of the task.
 */
export interface DeleteTaskPayload {
    userId: number;
    taskId: number;
}

/**
 * Checks if the task is valid.
 * @function isTaskValid
 * @param {TaskDto | null | undefined} task - The task to check.
 * @returns {task is TaskDto} - True if the task is valid, false otherwise.
 */
const isTaskValid = (task: TaskDto | null | undefined): task is TaskDto => {
    return !!(task && task.taskId && task.userId && task.name);
};

/**
 * Thunk for fetching a task.
 * @function fetchTask
 * @param {GetUserTaskPayload} payload - The payload containing userId and taskId.
 * @param {Object} thunkAPI - The thunk API object.
 * @param {Function} thunkAPI.rejectWithValue - Function to reject with a value.
 * @returns {Promise<TaskDto | string>} - The fetched task data or an error message.
 */
export const fetchTask = createAsyncThunk(
    "tasks/fetchTask",
    async ({ userId, taskId }: GetUserTaskPayload, { rejectWithValue }) => {
        logger.debug(`Fetching task (ID: ${ taskId }) for user (ID: ${ userId })`);
        const taskClient = new TaskClient(AxiosClient);

        try {
            const response: ResponseDto<TaskDto> = await taskClient.getTask(userId, taskId);
            const task = response?.data;

            if (!isTaskValid(task)) {
                logger.warn(`Invalid task data received for task (ID: ${ taskId })`);
                return rejectWithValue("Invalid task data received");
            }

            logger.info(`Task (ID: ${ taskId }) successfully fetched for user (ID: ${ userId })`);
            return task;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Error fetching task");
            logger.error(`Failed to fetch task (ID: ${ taskId }) for user (ID: ${ userId }): ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);

/**
 * Thunk for fetching all tasks for a user.
 * @function fetchTasks
 * @param {number} userId - The ID of the user.
 * @param {Object} thunkAPI - The thunk API object.
 * @param {Function} thunkAPI.rejectWithValue - Function to reject with a value.
 * @returns {Promise<TaskDto[] | string>} - The fetched tasks data or an error message.
 */
export const fetchTasks = createAsyncThunk(
    "tasks/fetchTasks",
    async (userId: number, { rejectWithValue }) => {
        logger.debug(`Fetching tasks for user (ID: ${ userId })`);
        const taskClient = new TaskClient(AxiosClient);

        try {
            const response = await taskClient.getTasks(userId);
            const tasks = response?.data;

            if (!Array.isArray(tasks)) {
                logger.warn(`Invalid task list received for user (ID: ${ userId })`);
                return rejectWithValue("Invalid task list received");
            }

            logger.info(`Tasks successfully fetched for user (ID: ${ userId })`);
            return tasks;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Error fetching tasks");
            logger.error(`Failed to fetch tasks for user (ID: ${ userId }): ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);

/**
 * Thunk for creating a task.
 * @function createTask
 * @param {CreateTaskPayload} payload - The payload containing userId and taskData.
 * @param {Object} thunkAPI - The thunk API object.
 * @param {Function} thunkAPI.rejectWithValue - Function to reject with a value.
 * @returns {Promise<TaskDto | string>} - The created task data or an error message.
 */
export const createTask = createAsyncThunk(
    "tasks/createTask",
    async ({ userId, taskData }: CreateTaskPayload, { rejectWithValue }) => {
        logger.debug(`Creating a task for user (ID: ${ userId })`);
        const taskClient = new TaskClient(AxiosClient);

        try {
            const response = await taskClient.createTask(userId, taskData);
            const task = response?.data;

            if (!isTaskValid(task)) {
                logger.warn(`Invalid task data received after creation for user (ID: ${ userId })`);
                return rejectWithValue("Invalid task data received after creation");
            }

            logger.info(`Task successfully created for user (ID: ${ userId })`);
            return task;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Error creating task");
            logger.error(`Failed to create task for user (ID: ${ userId }): ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);

/**
 * Thunk for updating a task.
 * @function updateTask
 * @param {UpdateTaskPayload} payload - The payload containing userId, taskId, and taskData.
 * @param {Object} thunkAPI - The thunk API object.
 * @param {Function} thunkAPI.rejectWithValue - Function to reject with a value.
 * @returns {Promise<TaskDto | string>} - The updated task data or an error message.
 */
export const updateTask = createAsyncThunk(
    "tasks/updateTask",
    async ({ userId, taskId, taskData }: UpdateTaskPayload, { rejectWithValue }) => {
        logger.debug(`Updating task (ID: ${ taskId }) for user (ID: ${ userId })`);
        const taskClient = new TaskClient(AxiosClient);

        try {
            const response = await taskClient.updateTask(userId, taskId, taskData);
            const task = response?.data;

            if (!isTaskValid(task)) {
                logger.warn(`Invalid task data received after update for task (ID: ${ taskId })`);
                return rejectWithValue("Invalid task data received after update");
            }

            logger.info(`Task (ID: ${ taskId }) successfully updated for user (ID: ${ userId })`);
            return task;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Error updating task");
            logger.error(`Failed to update task (ID: ${ taskId }) for user (ID: ${ userId }): ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);

/**
 * Thunk for deleting a task.
 * @function deleteTask
 * @param {DeleteTaskPayload} payload - The payload containing userId and taskId.
 * @param {Object} thunkAPI - The thunk API object.
 * @param {Function} thunkAPI.rejectWithValue - Function to reject with a value.
 * @returns {Promise<null | string>} - Null if the task was deleted successfully, or an error message.
 */
export const deleteTask = createAsyncThunk(
    "tasks/deleteTask",
    async ({ userId, taskId }: DeleteTaskPayload, { rejectWithValue }) => {
        logger.debug(`Deleting task (ID: ${ taskId }) for user (ID: ${ userId })`);
        const taskClient = new TaskClient(AxiosClient);

        try {
            const response = await taskClient.deleteTask(userId, taskId);
            logger.info(`Task (ID: ${ taskId }) successfully deleted for user (ID: ${ userId })`);
            return response?.data;
        } catch (error: unknown) {
            const errorMessage = parseErrorMessage(error, "Error deleting task");
            logger.error(`Failed to delete task (ID: ${ taskId }) for user (ID: ${ userId }): ${ errorMessage }`);
            return rejectWithValue(errorMessage);
        }
    }
);