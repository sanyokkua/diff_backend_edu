import { DeleteResult, Repository, UpdateResult } from "typeorm";
import { ITaskRepository }                        from "../api";
import { Task, User }                             from "../models";


export class TaskRepository implements ITaskRepository {
    private readonly repository: Repository<Task>;

    constructor(repository: Repository<Task>) {
        this.repository = repository;
    }

    async findByTaskID(taskId: number): Promise<Task | null> {
        return await this.repository.findOneBy({ id: taskId });
    }

    async findByUserAndName(user: User, name: string): Promise<Task | null> {
        return await this.repository.findOne({ where: { user: { id: user.id }, name: name } });
    }

    async findByUserAndTaskID(user: User, taskId: number): Promise<Task | null> {
        return await this.repository.findOne({ where: { user: { id: user.id }, id: taskId } });
    }

    async findAllByUser(user: User): Promise<Task[]> {
        return await this.repository.find({ where: { user: { id: user.id } } });
    }

    async createTask(task: Task): Promise<Task> {
        return await this.repository.save(task);
    }

    async updateTask(id: number, updateData: Partial<Task>): Promise<UpdateResult> {
        return await this.repository.update(id, updateData);
    }

    async deleteTask(task: Task): Promise<DeleteResult> {
        return await this.repository.delete({ id: task.id });
    }
}
