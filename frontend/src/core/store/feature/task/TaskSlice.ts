/**
 * Redux slice for managing task-related state in the application.
 * This slice handles fetching, creating, updating, and deleting tasks,
 * and updates the task state accordingly. It integrates with Redux Toolkit
 * and uses thunks for asynchronous task operations. Logging is performed
 * using a custom logger at various stages to trace the lifecycle of each action.
 *
 * @module taskSlice
 */

import { createSlice, PayloadAction }                                from "@reduxjs/toolkit";
import { TaskDto }                                                   from "../../../api";
import { LogLevel }                                                  from "../../../config";
import { createTask, deleteTask, fetchTask, fetchTasks, updateTask } from "../../thunks";


const logger = LogLevel.getLogger("TaskSlice");

/**
 * TaskState interface defines the structure of the task-related state
 * that is managed by this slice.
 *
 * @interface TaskState
 * @property {number} taskId - ID of the currently selected task.
 * @property {string} taskName - Name of the currently selected task.
 * @property {string} taskDescription - Description of the currently selected task.
 * @property {TaskDto[]} tasksList - List of all tasks retrieved from the server.
 * @property {boolean} isTaskLoading - Indicates whether a task-related action is currently loading.
 * @property {string} taskError - Error message for any task-related failures.
 */
export interface TaskState {
    taskId: number;
    taskName: string;
    taskDescription: string;
    tasksList: TaskDto[];
    isTaskLoading: boolean;
    taskError: string;
}

/**
 * The initial state of the task slice, containing default values for each field.
 *
 * @constant
 * @type {TaskState}
 */
const initialState: TaskState = {
    taskId: 0,
    taskName: "",
    taskDescription: "",
    tasksList: [],
    isTaskLoading: false,
    taskError: ""
};

export const taskSlice = createSlice(
    {
        name: "tasks",
        initialState,
        reducers: {
            /**
             * Sets the currently selected task's ID.
             *
             * @param {TaskState} state - The current state of the slice.
             * @param {PayloadAction<number>} action - The action payload containing the task ID.
             */
            setTaskId: (state: TaskState, action: PayloadAction<number>) => {
                logger.debug(`Setting current task ID: ${ action.payload }`);
                state.taskId = action.payload;
            },

            /**
             * Sets the currently selected task's name.
             *
             * @param {TaskState} state - The current state of the slice.
             * @param {PayloadAction<string>} action - The action payload containing the task name.
             */
            setTaskName: (state: TaskState, action: PayloadAction<string>) => {
                logger.debug(`Setting current task name: ${ action.payload }`);
                state.taskName = action.payload;
            },

            /**
             * Sets the description of the currently selected task.
             *
             * @param {TaskState} state - The current state of the slice.
             * @param {PayloadAction<string>} action - The action payload containing the task description.
             */
            setTaskDescription: (state: TaskState, action: PayloadAction<string>) => {
                logger.debug(`Setting current task description: ${ action.payload }`);
                state.taskDescription = action.payload;
            },

            /**
             * Sets the list of all tasks retrieved from the server.
             *
             * @param {TaskState} state - The current state of the slice.
             * @param {PayloadAction<TaskDto[]>} action - The action payload containing an array of tasks.
             */
            setTasksList: (state: TaskState, action: PayloadAction<TaskDto[]>) => {
                logger.debug(`Setting task list: ${ action.payload.length } tasks`);
                state.tasksList = action.payload;
            },

            /**
             * Sets the loading state for task-related actions.
             *
             * @param {TaskState} state - The current state of the slice.
             * @param {PayloadAction<boolean>} action - The action payload indicating whether a task action is loading.
             */
            setIsTaskLoading: (state: TaskState, action: PayloadAction<boolean>) => {
                logger.debug(`Setting loading state: ${ action.payload }`);
                state.isTaskLoading = action.payload;
            },

            /**
             * Sets the error message for any task-related action that fails.
             *
             * @param {TaskState} state - The current state of the slice.
             * @param {PayloadAction<string>} action - The action payload containing the error message.
             */
            setTaskError: (state: TaskState, action: PayloadAction<string>) => {
                logger.debug(`Setting error message: ${ action.payload }`);
                state.taskError = action.payload;
            }
        },
        extraReducers: (builder) => {
            builder
                // Handle pending state for fetching a single task
                .addCase(fetchTask.pending, (state: TaskState) => {
                    logger.debug("fetchTask request pending");
                    state.isTaskLoading = true;
                    state.taskError = "";
                })

                // Handle fulfilled state for fetching a single task
                .addCase(fetchTask.fulfilled, (state: TaskState, action: PayloadAction<TaskDto>) => {
                    logger.info("fetchTask request succeeded");
                    state.isTaskLoading = false;
                    state.taskError = "";

                    const task = action.payload;
                    state.taskId = task.taskId;
                    state.taskName = task.name;
                    state.taskDescription = task.description;
                })

                // Handle rejected state for fetching a single task
                .addCase(fetchTask.rejected, (state: TaskState, action) => {
                    logger.error("fetchTask request failed");
                    state.isTaskLoading = false;
                    state.taskError = action.payload as string || "Failed to fetch task.";
                })

                // Handle pending state for fetching all tasks
                .addCase(fetchTasks.pending, (state: TaskState) => {
                    logger.debug("fetchTasks request pending");
                    state.isTaskLoading = true;
                    state.taskError = "";
                    state.tasksList = [];
                })

                // Handle fulfilled state for fetching all tasks
                .addCase(fetchTasks.fulfilled, (state: TaskState, action: PayloadAction<TaskDto[]>) => {
                    logger.info("fetchTasks request succeeded");
                    state.isTaskLoading = false;
                    state.taskError = "";
                    state.tasksList = action.payload;
                })

                // Handle rejected state for fetching all tasks
                .addCase(fetchTasks.rejected, (state: TaskState, action) => {
                    logger.error("fetchTasks request failed");
                    state.isTaskLoading = false;
                    state.taskError = action.payload as string || "Failed to fetch tasks.";
                })

                // Handle pending state for creating a task
                .addCase(createTask.pending, (state: TaskState) => {
                    logger.debug("createTask request pending");
                    state.isTaskLoading = true;
                    state.taskError = "";
                })

                // Handle fulfilled state for creating a task
                .addCase(createTask.fulfilled, (state: TaskState, action: PayloadAction<TaskDto>) => {
                    logger.info("createTask request succeeded");
                    state.isTaskLoading = false;
                    state.taskError = "";

                    const task = action.payload;
                    state.taskId = task.taskId;
                    state.taskName = task.name;
                    state.taskDescription = task.description;
                })

                // Handle rejected state for creating a task
                .addCase(createTask.rejected, (state: TaskState, action) => {
                    logger.error("createTask request failed");
                    state.isTaskLoading = false;
                    state.taskError = action.payload as string || "Failed to create task.";
                })

                // Handle pending state for updating a task
                .addCase(updateTask.pending, (state: TaskState) => {
                    logger.debug("updateTask request pending");
                    state.isTaskLoading = true;
                    state.taskError = "";
                })

                // Handle fulfilled state for updating a task
                .addCase(updateTask.fulfilled, (state: TaskState, action: PayloadAction<TaskDto>) => {
                    logger.info("updateTask request succeeded");
                    state.isTaskLoading = false;
                    state.taskError = "";

                    const task = action.payload;
                    state.taskId = task.taskId;
                    state.taskName = task.name;
                    state.taskDescription = task.description;
                })

                // Handle rejected state for updating a task
                .addCase(updateTask.rejected, (state: TaskState, action) => {
                    logger.error("updateTask request failed");
                    state.isTaskLoading = false;
                    state.taskError = action.payload as string || "Failed to update task.";
                })

                // Handle pending state for deleting a task
                .addCase(deleteTask.pending, (state: TaskState) => {
                    logger.debug("deleteTask request pending");
                    state.isTaskLoading = true;
                    state.taskError = "";
                })

                // Handle fulfilled state for deleting a task
                .addCase(deleteTask.fulfilled, (state: TaskState) => {
                    logger.info("deleteTask request succeeded");
                    state.isTaskLoading = false;
                    state.taskError = "";

                    // Reset task details after deletion
                    state.taskId = 0;
                    state.taskName = "";
                    state.taskDescription = "";
                })

                // Handle rejected state for deleting a task
                .addCase(deleteTask.rejected, (state: TaskState, action) => {
                    logger.error("deleteTask request failed");
                    state.isTaskLoading = false;
                    state.taskError = action.payload as string || "Failed to delete task.";
                });
        }
    }
);

/**
 * User slice action creators.
 *
 * @property {function} setTaskId - Action to set the TaskId.
 * @property {function} setTaskName - Action to set the TaskName.
 * @property {function} setTaskDescription - Action to set the TaskDescription.
 * @property {function} setTasksList - Action to set the TasksList.
 * @property {function} setIsTaskLoading - Action to set the IsTaskLoading.
 * @property {function} setTaskError - Action to set an TaskError.
 */
export const {
    setTaskId,
    setTaskName,
    setTaskDescription,
    setTasksList,
    setIsTaskLoading,
    setTaskError
} = taskSlice.actions;

// Export the task slice reducer
export default taskSlice.reducer;
