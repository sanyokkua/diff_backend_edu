import { DeleteResult, Repository, UpdateResult } from "typeorm";
import { ITaskRepository }                        from "../api";
import { IllegalArgumentError }                   from "../error";
import { Task, User }                             from "../model";


export class TaskRepository implements ITaskRepository {
    private readonly repository: Repository<Task>;

    constructor(repository: Repository<Task>) {
        this.repository = repository;
    }

    async findByTaskID(taskId: number): Promise<Task | null> {
        if (taskId <= 0) {
            throw new IllegalArgumentError("ID is not valid");
        }
        return await this.findTask({ id: taskId });
    }

    async findByUserAndName(user: User, name: string): Promise<Task | null> {
        if (!user || !name.trim()) {
            throw new IllegalArgumentError("Invalid params");
        }
        return await this.findTask({ user: user, name: name });
    }

    async findByUserAndTaskID(user: User, taskId: number): Promise<Task | null> {
        if (!user || taskId <= 0) {
            throw new IllegalArgumentError("Invalid params");
        }
        return await this.findTask({ user: user, id: taskId });
    }

    async findAllByUser(user: User): Promise<Task[]> {
        if (!user) {
            throw new IllegalArgumentError("User is null");
        }
        try {
            const tasks = await this.repository.find({ where: { user: { id: user.id } } });
            return tasks;
        } catch (error) {
            throw error;
        }
    }

    async createTask(task: Task): Promise<Task> {
        if (!task) {
            throw new IllegalArgumentError("Task model is null");
        }
        try {
            const newTask = await this.repository.save(task);
            return newTask;
        } catch (error) {
            throw error;
        }
    }

    async updateTask(id: number, updateData: Partial<Task>): Promise<UpdateResult> {
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

    async deleteTask(task: Task): Promise<DeleteResult> {
        if (!task) {
            throw new IllegalArgumentError("Task model is nil");
        }
        try {
            const result = await this.repository.delete({ id: task.id });
            return result;
        } catch (error) {
            throw error;
        }
    }

    // Private method to encapsulate task retrieval logic.
    private async findTask(criteria: Partial<Task>): Promise<Task | null> {
        try {
            const task = await this.repository.findOne({ where: criteria });
            // if (task) {
            //     Logger.info(`Successfully retrieved task with id: ${task.id}`);
            // } else {
            //     Logger.warn(`Task not found for criteria: ${JSON.stringify(criteria)}`);
            // }
            return task;
        } catch (error) {
            throw error;
        }
    }
}
