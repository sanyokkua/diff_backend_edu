export interface TaskCreationDTO {
    name: string;
    description: string;
}

export interface TaskDto {
    userId: number;
    taskId: number;
    name: string;
    description: string;
}

export interface TaskUpdateDTO {
    name: string;
    description: string;
}
