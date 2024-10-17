import { DeleteResult, Repository, UpdateResult } from "typeorm";
import { IUserRepository }                        from "../api";
import { User }                                   from "../models";


export class UserRepository implements IUserRepository {
    private readonly repository: Repository<User>;

    constructor(repository: Repository<User>) {
        this.repository = repository;
    }

    async createUser(user: User): Promise<User> {
        return await this.repository.save(user);
    }

    async deleteUser(id: number): Promise<DeleteResult> {
        return await this.repository.delete(id);
    }

    async getUserByEmail(email: string): Promise<User | null> {
        return await this.repository.findOneBy({ email });
    }

    async getUserByID(id: number): Promise<User | null> {
        return await this.repository.findOneBy({ id });
    }

    async updateUser(id: number, updateData: Partial<User>): Promise<UpdateResult> {
        return await this.repository.update(id, updateData);
    }
}