import { DeleteResult, UpdateResult }                               from "typeorm";
import { UserCreationDTO, UserDeletionDTO, UserDTO, UserUpdateDTO } from "../dto";
import { User }                                                     from "../models";


export interface IUserRepository {
    createUser(user: User): Promise<User>;

    getUserByID(id: number): Promise<User | null>;

    updateUser(id: number, updateData: Partial<User>): Promise<UpdateResult>;

    deleteUser(id: number): Promise<DeleteResult>;

    getUserByEmail(email: string): Promise<User | null>;
}

export interface IUserService {
    create(userCreationDTO: UserCreationDTO): Promise<UserDTO>;

    updatePassword(userId: number, userUpdateDTO: UserUpdateDTO): Promise<UserDTO>;

    delete(userId: number, userDeletionDTO: UserDeletionDTO): Promise<void>;
}
