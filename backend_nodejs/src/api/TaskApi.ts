import { DeleteResult, UpdateResult }              from "typeorm";
import { TaskCreationDTO, TaskDTO, TaskUpdateDTO } from "../dto";
import { Task, User }                              from "../models";


export interface ITaskRepository {
    findByTaskID(taskId: number): Promise<Task | null>;

    findByUserAndName(user: User, name: string): Promise<Task | null>;

    findByUserAndTaskID(user: User, taskId: number): Promise<Task | null>;

    findAllByUser(user: User): Promise<Task[]>;

    createTask(task: Task): Promise<Task>;

    updateTask(id: number, updateData: Partial<Task>): Promise<UpdateResult>;

    deleteTask(task: Task): Promise<DeleteResult>;
}

export interface ITaskService {
    createTask(userId: number, taskCreationDTO: TaskCreationDTO): Promise<TaskDTO>;

    updateTask(userId: number, taskID: number, taskUpdateDTO: TaskUpdateDTO): Promise<TaskDTO>;

    deleteTask(userId: number, taskID: number): Promise<void>;

    getAllTasksForUser(userId: number): Promise<TaskDTO[]>;

    getTaskByUserIDAndTaskID(userId: number, taskID: number): Promise<TaskDTO>;
}
