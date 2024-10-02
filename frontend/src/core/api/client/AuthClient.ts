import { AxiosInstance }                 from "axios";
import { UserCreationDTO, UserLoginDto } from "../dto/Auth";
import { ResponseDto }                   from "../dto/Common";
import { UserDto }                       from "../dto/User";
import { handleError, handleResponse }   from "./Utils";


const BASE_URL = "/auth";

class AuthClient {
    private readonly axiosClient: AxiosInstance;

    constructor(axiosClient: AxiosInstance) {
        this.axiosClient = axiosClient;
    }

    async loginUser(request: UserLoginDto): Promise<ResponseDto<UserDto>> {
        try {
            const url: string = `${ BASE_URL }/login`;
            const response = await this.axiosClient.post<ResponseDto<UserDto>>(url, request);
            return handleResponse(response);
        } catch (error) {
            return handleError(error);
        }
    }

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