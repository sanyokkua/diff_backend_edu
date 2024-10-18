export type TaskDto = {
    taskId: number
    name: string
    description: string
    userId: number
}

export type TaskCreationDto = {
    name: string
    description: string
}

export type TaskUpdateDto = {
    name: string
    description: string
}
