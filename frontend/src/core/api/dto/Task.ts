/**
 * Data Transfer Object for task creation.
 *
 * @property {string} name - The name of the task.
 * @property {string} description - The description of the task.
 */
export interface TaskCreationDTO {
    name: string;
    description: string;
}

/**
 * Data Transfer Object for task data.
 *
 * @property {number} userId - The ID of the user to whom the task belongs.
 * @property {number} taskId - The unique identifier of the task.
 * @property {string} name - The name of the task.
 * @property {string} description - The description of the task.
 */
export interface TaskDto {
    userId: number;
    taskId: number;
    name: string;
    description: string;
}

/**
 * Data Transfer Object for task update.
 *
 * @property {string} name - The updated name of the task.
 * @property {string} description - The updated description of the task.
 */
export interface TaskUpdateDTO {
    name: string;
    description: string;
}
