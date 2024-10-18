import { DeleteResult, UpdateResult }              from "typeorm";
import { TaskCreationDto, TaskDto, TaskUpdateDto } from "../dto";
import { Task, User }                              from "../model";


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
    createTask(userId: number, taskCreationDTO: TaskCreationDto): Promise<TaskDto>;

    updateTask(userId: number, taskID: number, taskUpdateDTO: TaskUpdateDto): Promise<TaskDto>;

    deleteTask(userId: number, taskID: number): Promise<void>;

    getAllTasksForUser(userId: number): Promise<TaskDto[]>;

    getTaskByUserIDAndTaskID(userId: number, taskID: number): Promise<TaskDto>;
}
