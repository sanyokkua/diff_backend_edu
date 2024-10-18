import { ITaskRepository, ITaskService, IUserRepository }                      from "../api";
import { TaskCreationDto, TaskDto, TaskUpdateDto }                             from "../dto";
import { IllegalArgumentError, TaskNotFoundError }                             from "../error";
import { Task }                                                                from "../model";
import { checkTaskExistsForUser, validateTaskCreation, validateTaskUpdateDTO } from "../utils";


export class TaskService implements ITaskService {
    private readonly taskRepository: ITaskRepository;
    private readonly userRepository: IUserRepository;

    constructor(taskRepository: ITaskRepository, userRepository: IUserRepository) {
        this.taskRepository = taskRepository;
        this.userRepository = userRepository;
    }

    async createTask(userId: number, taskCreationDTO: TaskCreationDto): Promise<TaskDto> {
        validateTaskCreation(taskCreationDTO);
        const userByID = await this.userRepository.getUserByID(userId);
        if (!userByID) {
            throw new IllegalArgumentError("User not found");
        }
        await checkTaskExistsForUser(this.taskRepository, userByID, taskCreationDTO.name);

        const task: Task = {
            name: taskCreationDTO.name,
            description: taskCreationDTO.description,
            user: userByID
        };

        const createdTask = await this.taskRepository.createTask(task);

        return {
            taskId: createdTask.id ?? -1,
            name: createdTask.name,
            description: createdTask.description,
            userId: createdTask.user.id ?? -1
        };
    }

    async updateTask(userId: number, taskID: number, taskUpdateDTO: TaskUpdateDto): Promise<TaskDto> {
        validateTaskUpdateDTO(taskUpdateDTO);
        const userByID = await this.userRepository.getUserByID(userId);
        if (!userByID) {
            throw new IllegalArgumentError("User not found");
        }

        const taskFromRepo = await this.taskRepository.findByUserAndTaskID(userByID, taskID);
        if (!taskFromRepo) {
            throw new TaskNotFoundError("Task is not found");
        }

        taskFromRepo.name = taskUpdateDTO.name;
        taskFromRepo.description = taskUpdateDTO.description;

        const updateResult = await this.taskRepository.updateTask(taskID, taskFromRepo);

        return {
            taskId: taskFromRepo.id ?? -1,
            name: taskFromRepo.name,
            description: taskFromRepo.description,
            userId: userByID.id ?? -1
        };
    }

    async deleteTask(userId: number, taskID: number): Promise<void> {
        const userByID = await this.userRepository.getUserByID(userId);
        if (!userByID) {
            throw new IllegalArgumentError("User not found");
        }

        const foundTask = await this.taskRepository.findByUserAndTaskID(userByID, taskID);
        if (!foundTask) {
            throw new TaskNotFoundError("Task is not found");
        }

        const deleteResult = await this.taskRepository.deleteTask(foundTask);
    }

    async getAllTasksForUser(userId: number): Promise<TaskDto[]> {
        const userByID = await this.userRepository.getUserByID(userId);
        if (!userByID) {
            throw new IllegalArgumentError("User not found");
        }

        const tasks: Task[] = await this.taskRepository.findAllByUser(userByID);
        return tasks.map(t => {
            return {
                taskId: t.id,
                userId: userByID.id,
                name: t.name,
                description: t.description
            } as TaskDto;
        });
    }

    async getTaskByUserIDAndTaskID(userId: number, taskID: number): Promise<TaskDto> {
        const userByID = await this.userRepository.getUserByID(userId);
        if (!userByID) {
            throw new IllegalArgumentError("User not found");
        }

        const found = await this.taskRepository.findByUserAndTaskID(userByID, taskID);
        if (!found) {
            throw new TaskNotFoundError("Task is not found");
        }

        return {
            taskId: found.id,
            userId: found.user.id,
            name: found.name,
            description: found.description
        } as TaskDto;
    }

}