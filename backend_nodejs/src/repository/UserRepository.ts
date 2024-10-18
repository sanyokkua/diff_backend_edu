import { DeleteResult, Repository, UpdateResult }        from "typeorm";
import { IUserRepository }                               from "../api";
import { IllegalArgumentError, InvalidEmailFormatError } from "../error";
import { User }                                          from "../model";


export class UserRepository implements IUserRepository {
    private readonly repository: Repository<User>;

    constructor(repository: Repository<User>) {
        this.repository = repository;
    }

    async createUser(user: User): Promise<User> {
        if (!user) {
            throw new IllegalArgumentError("Passed User model is nil");
        }
        try {
            const newUser = await this.repository.save(user);
            return newUser;
        } catch (error) {
            throw error;
        }
    }

    async getUserByID(id: number): Promise<User | null> {
        if (id <= 0) {
            throw new IllegalArgumentError("Passed id is invalid");
        }
        try {
            const user = await this.repository.findOneBy({ id });
            // if (user) {
            //     Logger.debug(`User retrieved successfully with ID: ${ id }`);
            // } else {
            //     Logger.warn(`No user found with ID: ${ id }`);
            // }
            return user;
        } catch (error) {
            throw error;
        }
    }

    async updateUser(id: number, updateData: Partial<User>): Promise<UpdateResult> {
        if (id <= 0 || !updateData) {
            throw new IllegalArgumentError("Invalid params");
        }
        try {
            const result = await this.repository.update(id, updateData);
            return result;
        } catch (error) {
            throw error;
        }
    }

    async deleteUser(id: number): Promise<DeleteResult> {
        if (id <= 0) {
            throw new IllegalArgumentError("Passed id is invalid");
        }
        try {
            const result = await this.repository.delete(id);
            return result;
        } catch (error) {
            throw error;
        }
    }

    async getUserByEmail(email: string): Promise<User | null> {
        if (!email || !email.trim()) {
            throw new InvalidEmailFormatError("Email can't be an empty string");
        }
        try {
            const user = await this.repository.findOneBy({ email });
            // if (user) {
            //     Logger.debug(`User retrieved successfully with email: ${ email }`);
            // } else {
            //     Logger.warn(`No user found with email: ${ email }`);
            // }
            return user;
        } catch (error) {
            throw error;
        }
    }
}
