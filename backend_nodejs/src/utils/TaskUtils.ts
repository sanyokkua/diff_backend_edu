import { ITaskRepository }                              from "../api";
import { TaskCreationDto, TaskUpdateDto }               from "../dto";
import { IllegalArgumentError, TaskAlreadyExistsError } from "../error";
import { User }                                         from "../model";


export async function checkTaskExistsForUser(taskRepo: ITaskRepository, user: User, taskName: string): Promise<void> {
    try {
        const task = await taskRepo.findByUserAndName(user, taskName);
        if (task) {
            throw new TaskAlreadyExistsError(`Task with the name '${ taskName }' already exists for the user`);
        }
    } catch (err) {
        if ((err as any)?.message !== "RecordNotFound") {
            throw err;
        }
    }
}

export function validateTaskCreation(taskCreationDTO: TaskCreationDto | null): void {
    if (taskCreationDTO === null) {
        throw new IllegalArgumentError("TaskCreationDto is nil");
    }
    if (taskCreationDTO.name === "") {
        throw new IllegalArgumentError("Task name cannot be empty");
    }
    if (taskCreationDTO.description === "") {
        throw new IllegalArgumentError("Task description cannot be empty");
    }
}

export function validateTaskUpdateDTO(taskUpdateDTO: TaskUpdateDto | null): void {
    if (taskUpdateDTO === null) {
        throw new IllegalArgumentError("TaskUpdateDto is nil");
    }
    if (taskUpdateDTO.name === "") {
        throw new IllegalArgumentError("Task name cannot be empty");
    }
    if (taskUpdateDTO.description === "") {
        throw new IllegalArgumentError("Task description cannot be empty");
    }
}
