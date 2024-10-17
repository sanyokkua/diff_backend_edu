import { ITaskRepository, ITaskService, IUserRepository } from "../api";
import { TaskCreationDTO, TaskDTO, TaskUpdateDTO }        from "../dto";


export class TaskService implements ITaskService {
    private readonly taskRepository: ITaskRepository;
    private readonly userRepository: IUserRepository;

    constructor(taskRepository: ITaskRepository, userRepository: IUserRepository) {
        this.taskRepository = taskRepository;
        this.userRepository = userRepository;
    }

    createTask(userId: number, taskCreationDTO: TaskCreationDTO): Promise<TaskDTO> {
        throw new Error("Method not implemented.");
    }

    updateTask(userId: number, taskID: number, taskUpdateDTO: TaskUpdateDTO): Promise<TaskDTO> {
        throw new Error("Method not implemented.");
    }

    deleteTask(userId: number, taskID: number): Promise<void> {
        throw new Error("Method not implemented.");
    }

    getAllTasksForUser(userId: number): Promise<TaskDTO[]> {
        throw new Error("Method not implemented.");
    }

    getTaskByUserIDAndTaskID(userId: number, taskID: number): Promise<TaskDTO> {
        throw new Error("Method not implemented.");
    }

}