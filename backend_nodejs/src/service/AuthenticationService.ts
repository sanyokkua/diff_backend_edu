import { IAuthenticationService, IJwtService, IPasswordEncoder, IUserRepository, IUserService } from "../api";
import { UserCreationDto, UserDto, UserLoginDto }                                               from "../dto";
import { IllegalArgumentError, InvalidPasswordError }                                           from "../error";
import { validateUserCreationDTO, validateUserLoginDto }                                        from "../utils";


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

    async loginUser(dto: UserLoginDto): Promise<UserDto> {
        validateUserLoginDto(dto);
        const userByEmail = await this.userRepository.getUserByEmail(dto.email);
        if (!userByEmail) {
            throw new IllegalArgumentError("User Not Found");
        }

        const matches: boolean = this.passwordEncoder.matches(dto.password, userByEmail.passwordHash);
        if (!matches) {
            throw new InvalidPasswordError("Invalid credentials");
        }

        const jwtToken: string = await this.jwtService.generateJwtToken(userByEmail.email);

        return {
            userId: userByEmail.id ?? -1,
            email: userByEmail.email,
            jwtToken: jwtToken
        };
    }

    async registerUser(dto: UserCreationDto): Promise<UserDto> {
        validateUserCreationDTO(dto);

        const newUser = await this.userService.create(dto);
        const jwtToken: string = await this.jwtService.generateJwtToken(newUser.email);

        return {
            userId: newUser.userId,
            email: newUser.email,
            jwtToken: jwtToken
        };
    }

}