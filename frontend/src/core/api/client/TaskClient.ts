import { AxiosInstance }                           from "axios";
import { ResponseDto }                             from "../dto/Common";
import { TaskCreationDTO, TaskDto, TaskUpdateDTO } from "../dto/Task";
import { handleError, handleResponse }             from "./Utils";


const BASE_URL = "/users";

class TaskClient {
    private readonly axiosClient: AxiosInstance;

    constructor(axiosClient: AxiosInstance) {
        this.axiosClient = axiosClient;
    }

    async getTask(userId: number, taskId: number): Promise<ResponseDto<TaskDto>> {
        try {
            const url: string = `${ BASE_URL }/${ userId }/tasks/${ taskId }`;
            const response = await this.axiosClient.get<ResponseDto<TaskDto>>(url);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    }

    async getTasks(userId: number): Promise<ResponseDto<TaskDto[]>> {
        try {
            const url: string = `${ BASE_URL }/${ userId }/tasks/`;
            const response = await this.axiosClient.get<ResponseDto<TaskDto[]>>(url);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    }

    async createTask(userId: number, request: TaskCreationDTO): Promise<ResponseDto<TaskDto>> {
        try {
            const url: string = `${ BASE_URL }/${ userId }/tasks/`;
            const response = await this.axiosClient.post<ResponseDto<TaskDto>>(url, request);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    }

    async updateTask(userId: number, taskId: number, taskUpdateDTO: TaskUpdateDTO): Promise<ResponseDto<TaskDto>> {
        try {
            const url: string = `${ BASE_URL }/${ userId }/tasks/${ taskId }`;
            const response = await this.axiosClient.put<ResponseDto<TaskDto>>(url, taskUpdateDTO);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    }

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