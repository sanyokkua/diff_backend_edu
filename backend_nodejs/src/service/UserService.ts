import { IPasswordEncoder, IUserRepository, IUserService } from "../api";
import {
    UserCreationDto,
    UserDeletionDto,
    UserDto,
    UserUpdateDto
}                                                          from "../dto";
import {
    IllegalArgumentError
}                                                          from "../error";
import {
    User
}                                                          from "../model";
import {
    validatePasswordUpdate,
    validateUserCreationDTO,
    validateUserDeletionDTO,
    validateUserUpdateDTO
}                                                          from "../utils";


export class UserService implements IUserService {
    private readonly userRepository: IUserRepository;
    private readonly passwordEncoder: IPasswordEncoder;

    constructor(userRepository: IUserRepository, passwordEncoder: IPasswordEncoder) {
        this.userRepository = userRepository;
        this.passwordEncoder = passwordEncoder;
    }

    async create(userCreationDTO: UserCreationDto): Promise<UserDto> {
        validateUserCreationDTO(userCreationDTO);

        const encodedPassword: string = this.passwordEncoder.encode(userCreationDTO.password);
        const user: User = {
            id: 0,
            email: userCreationDTO.email,
            passwordHash: encodedPassword
        };
        const created: User = await this.userRepository.createUser(user);

        return {
            userId: created.id ?? -1,
            email: created.email,
            jwtToken: ""
        };
    }

    async updatePassword(userId: number, userUpdateDTO: UserUpdateDto): Promise<UserDto> {
        validateUserUpdateDTO(userUpdateDTO);
        const userFromDb = await this.userRepository.getUserByID(userId);
        if (!userFromDb) {
            throw new IllegalArgumentError("User not found");
        }

        validatePasswordUpdate(userUpdateDTO, userFromDb, this.passwordEncoder);

        userFromDb.passwordHash = this.passwordEncoder.encode(userUpdateDTO.newPassword);
        const updateResult = await this.userRepository.updateUser(userId, userFromDb);

        return {
            userId: userFromDb.id ?? -1,
            email: userFromDb.email,
            jwtToken: ""
        };
    }

    async delete(userId: number, userDeletionDTO: UserDeletionDto): Promise<void> {
        validateUserDeletionDTO(userDeletionDTO);

        const userFromDb = await this.userRepository.getUserByID(userId);
        if (!userFromDb) {
            throw new IllegalArgumentError("User not found");
        }

        const deleteResult = await this.userRepository.deleteUser(userId);
    }

}