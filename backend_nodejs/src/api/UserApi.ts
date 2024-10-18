import { DeleteResult, UpdateResult }                               from "typeorm";
import { UserCreationDto, UserDeletionDto, UserDto, UserUpdateDto } from "../dto";
import { User }                                                     from "../model";


export interface IUserRepository {
    createUser(user: User): Promise<User>;

    getUserByID(id: number): Promise<User | null>;

    updateUser(id: number, updateData: Partial<User>): Promise<UpdateResult>;

    deleteUser(id: number): Promise<DeleteResult>;

    getUserByEmail(email: string): Promise<User | null>;
}

export interface IUserService {
    create(userCreationDTO: UserCreationDto): Promise<UserDto>;

    updatePassword(userId: number, userUpdateDTO: UserUpdateDto): Promise<UserDto>;

    delete(userId: number, userDeletionDTO: UserDeletionDto): Promise<void>;
}
