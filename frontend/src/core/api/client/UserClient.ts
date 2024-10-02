import { AxiosInstance }                           from "axios";
import { ResponseDto }                             from "../dto/Common";
import { UserDeletionDTO, UserDto, UserUpdateDTO } from "../dto/User";
import { handleError, handleResponse }             from "./Utils";


const BASE_URL = "/users";

class UserClient {
    private readonly axiosClient: AxiosInstance;

    constructor(axiosClient: AxiosInstance) {
        this.axiosClient = axiosClient;
    }

    async getUser(userId: number): Promise<ResponseDto<UserDto>> {
        try {
            const url: string = `${ BASE_URL }/${ userId }`;
            const response = await this.axiosClient.get<ResponseDto<UserDto>>(url);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    }

    async updateUser(userId: number, userUpdateDTO: UserUpdateDTO): Promise<ResponseDto<UserDto>> {
        try {
            const url: string = `${ BASE_URL }/${ userId }/password`;
            const response = await this.axiosClient.put<ResponseDto<UserDto>>(url, userUpdateDTO);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    }

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