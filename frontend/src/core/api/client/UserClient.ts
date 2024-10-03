import { AxiosInstance }                           from "axios";
import { ResponseDto }                             from "../dto/Common";
import { UserDeletionDTO, UserDto, UserUpdateDTO } from "../dto/User";
import { handleError, handleResponse }             from "./Utils";


const BASE_URL = "/users";

/**
 * UserClient class provides methods to interact with the user management API.
 * It includes methods for retrieving, updating, and deleting user data.
 */
class UserClient {
    private readonly axiosClient: AxiosInstance;

    /**
     * Creates an instance of UserClient.
     * @param {AxiosInstance} axiosClient - The Axios instance to be used for HTTP requests.
     */
    constructor(axiosClient: AxiosInstance) {
        this.axiosClient = axiosClient;
    }

    /**
     * Retrieves a specific user by their ID.
     * @param {number} userId - The ID of the user.
     * @returns {Promise<ResponseDto<UserDto>>} - A promise that resolves to the response data transfer object containing user data.
     * @throws Will throw an error if the request fails.
     */
    async getUser(userId: number): Promise<ResponseDto<UserDto>> {
        try {
            const url: string = `${ BASE_URL }/${ userId }`;
            const response = await this.axiosClient.get<ResponseDto<UserDto>>(url);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    }

    /**
     * Updates a user's password.
     * @param {number} userId - The ID of the user.
     * @param {UserUpdateDTO} userUpdateDTO - The user update request data transfer object containing updated user details.
     * @returns {Promise<ResponseDto<UserDto>>} - A promise that resolves to the response data transfer object containing the updated user data.
     * @throws Will throw an error if the request fails.
     */
    async updateUser(userId: number, userUpdateDTO: UserUpdateDTO): Promise<ResponseDto<UserDto>> {
        try {
            const url: string = `${ BASE_URL }/${ userId }/password`;
            const response = await this.axiosClient.put<ResponseDto<UserDto>>(url, userUpdateDTO);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    }

    /**
     * Deletes a specific user.
     * @param {number} userId - The ID of the user.
     * @param {UserDeletionDTO} userDeletionDTO - The user deletion request data transfer object containing user deletion details.
     * @returns {Promise<ResponseDto<void>>} - A promise that resolves to the response data transfer object indicating the result of the deletion.
     * @throws Will throw an error if the request fails.
     */
    async deleteUser(userId: number, userDeletionDTO: UserDeletionDTO): Promise<ResponseDto<void>> {
        try {
            const url: string = `${ BASE_URL }/${ userId }/delete`;
            const response = await this.axiosClient.post<ResponseDto<void>>(url, userDeletionDTO);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    }
}

export default UserClient;
