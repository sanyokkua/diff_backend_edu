import { AxiosInstance }                 from "axios";
import { UserCreationDTO, UserLoginDto } from "../dto/Auth";
import { ResponseDto }                   from "../dto/Common";
import { UserDto }                       from "../dto/User";
import { handleError, handleResponse }   from "./Utils";


const BASE_URL = "/auth";

/**
 * AuthClient class provides methods to interact with the authentication API.
 * It includes methods for user login and registration.
 */
class AuthClient {
    private readonly axiosClient: AxiosInstance;

    /**
     * Creates an instance of AuthClient.
     * @param {AxiosInstance} axiosClient - The Axios instance to be used for HTTP requests.
     */
    constructor(axiosClient: AxiosInstance) {
        this.axiosClient = axiosClient;
    }

    /**
     * Logs in a user with the provided credentials.
     * @param {UserLoginDto} request - The login request data transfer object containing user credentials.
     * @returns {Promise<ResponseDto<UserDto>>} - A promise that resolves to the response data transfer object containing user data.
     * @throws Will throw an error if the login request fails.
     */
    async loginUser(request: UserLoginDto): Promise<ResponseDto<UserDto>> {
        try {
            const url: string = `${ BASE_URL }/login`;
            const response = await this.axiosClient.post<ResponseDto<UserDto>>(url, request);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    }

    /**
     * Registers a new user with the provided details.
     * @param {UserCreationDTO} request - The registration request data transfer object containing user details.
     * @returns {Promise<ResponseDto<UserDto>>} - A promise that resolves to the response data transfer object containing user data.
     * @throws Will throw an error if the registration request fails.
     */
    async registerUser(request: UserCreationDTO): Promise<ResponseDto<UserDto>> {
        try {
            const url: string = `${ BASE_URL }/register`;
            const response = await this.axiosClient.post<ResponseDto<UserDto>>(url, request);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    }
}

export default AuthClient;
