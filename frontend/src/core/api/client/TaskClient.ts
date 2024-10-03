import { AxiosInstance }                           from "axios";
import { ResponseDto }                             from "../dto/Common";
import { TaskCreationDTO, TaskDto, TaskUpdateDTO } from "../dto/Task";
import { handleError, handleResponse }             from "./Utils";


const BASE_URL = "/users";

/**
 * TaskClient class provides methods to interact with the task management API.
 * It includes methods for retrieving, creating, updating, and deleting tasks.
 */
class TaskClient {
    private readonly axiosClient: AxiosInstance;

    /**
     * Creates an instance of TaskClient.
     * @param {AxiosInstance} axiosClient - The Axios instance to be used for HTTP requests.
     */
    constructor(axiosClient: AxiosInstance) {
        this.axiosClient = axiosClient;
    }

    /**
     * Retrieves a specific task for a user.
     * @param {number} userId - The ID of the user.
     * @param {number} taskId - The ID of the task.
     * @returns {Promise<ResponseDto<TaskDto>>} - A promise that resolves to the response data transfer object containing task data.
     * @throws Will throw an error if the request fails.
     */
    async getTask(userId: number, taskId: number): Promise<ResponseDto<TaskDto>> {
        try {
            const url: string = `${ BASE_URL }/${ userId }/tasks/${ taskId }`;
            const response = await this.axiosClient.get<ResponseDto<TaskDto>>(url);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    }

    /**
     * Retrieves all tasks for a user.
     * @param {number} userId - The ID of the user.
     * @returns {Promise<ResponseDto<TaskDto[]>>} - A promise that resolves to the response data transfer object containing an array of tasks.
     * @throws Will throw an error if the request fails.
     */
    async getTasks(userId: number): Promise<ResponseDto<TaskDto[]>> {
        try {
            const url: string = `${ BASE_URL }/${ userId }/tasks/`;
            const response = await this.axiosClient.get<ResponseDto<TaskDto[]>>(url);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    }

    /**
     * Creates a new task for a user.
     * @param {number} userId - The ID of the user.
     * @param {TaskCreationDTO} request - The task creation request data transfer object containing task details.
     * @returns {Promise<ResponseDto<TaskDto>>} - A promise that resolves to the response data transfer object containing the created task data.
     * @throws Will throw an error if the request fails.
     */
    async createTask(userId: number, request: TaskCreationDTO): Promise<ResponseDto<TaskDto>> {
        try {
            const url: string = `${ BASE_URL }/${ userId }/tasks/`;
            const response = await this.axiosClient.post<ResponseDto<TaskDto>>(url, request);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    }

    /**
     * Updates an existing task for a user.
     * @param {number} userId - The ID of the user.
     * @param {number} taskId - The ID of the task.
     * @param {TaskUpdateDTO} taskUpdateDTO - The task update request data transfer object containing updated task details.
     * @returns {Promise<ResponseDto<TaskDto>>} - A promise that resolves to the response data transfer object containing the updated task data.
     * @throws Will throw an error if the request fails.
     */
    async updateTask(userId: number, taskId: number, taskUpdateDTO: TaskUpdateDTO): Promise<ResponseDto<TaskDto>> {
        try {
            const url: string = `${ BASE_URL }/${ userId }/tasks/${ taskId }`;
            const response = await this.axiosClient.put<ResponseDto<TaskDto>>(url, taskUpdateDTO);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    }

    /**
     * Deletes a specific task for a user.
     * @param {number} userId - The ID of the user.
     * @param {number} taskId - The ID of the task.
     * @returns {Promise<ResponseDto<void>>} - A promise that resolves to the response data transfer object indicating the result of the deletion.
     * @throws Will throw an error if the request fails.
     */
    async deleteTask(userId: number, taskId: number): Promise<ResponseDto<void>> {
        try {
            const url: string = `${ BASE_URL }/${ userId }/tasks/${ taskId }`;
            const response = await this.axiosClient.delete<ResponseDto<void>>(url);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    }
}

export default TaskClient;
