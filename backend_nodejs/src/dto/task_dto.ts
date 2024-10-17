export type TaskDTO = {
    taskId: number
    name: string
    description: string
    userId: number
}

export type TaskCreationDTO = {
    name: string
    description: string
}

export type TaskUpdateDTO = {
    name: string
    description: string
}
