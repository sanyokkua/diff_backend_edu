import { IAuthenticationService, IJwtService, IPasswordEncoder, IUserRepository, IUserService } from "../api";
import { UserCreationDTO, UserDTO, UserLoginDTO }                                               from "../dto";


export class AuthenticationService implements IAuthenticationService {
    private readonly userService: IUserService;
    private readonly userRepository: IUserRepository;
    private readonly jwtService: IJwtService;
    private readonly passwordEncoder: IPasswordEncoder;

    constructor(userService: IUserService, userRepository: IUserRepository, jwtService: IJwtService, passwordEncoder: IPasswordEncoder) {
        this.userService = userService;
        this.userRepository = userRepository;
        this.jwtService = jwtService;
        this.passwordEncoder = passwordEncoder;
    }

    loginUser(dto: UserLoginDTO): Promise<UserDTO> {
        throw new Error("Method not implemented.");
    }

    registerUser(dto: UserCreationDTO): Promise<UserDTO> {
        throw new Error("Method not implemented.");
    }

}