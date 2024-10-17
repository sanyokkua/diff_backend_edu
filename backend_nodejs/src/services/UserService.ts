import { IPasswordEncoder, IUserRepository, IUserService }          from "../api";
import { UserCreationDTO, UserDeletionDTO, UserDTO, UserUpdateDTO } from "../dto";


export class UserService implements IUserService {
    private readonly userRepository: IUserRepository;
    private readonly passwordEncoder: IPasswordEncoder;

    constructor(userRepository: IUserRepository, passwordEncoder: IPasswordEncoder) {
        this.userRepository = userRepository;
        this.passwordEncoder = passwordEncoder;
    }

    create(userCreationDTO: UserCreationDTO): Promise<UserDTO> {
        throw new Error("Method not implemented.");
    }

    updatePassword(userId: number, userUpdateDTO: UserUpdateDTO): Promise<UserDTO> {
        throw new Error("Method not implemented.");
    }

    delete(userId: number, userDeletionDTO: UserDeletionDTO): Promise<void> {
        throw new Error("Method not implemented.");
    }

}