import { Request, Response } from "express";
import { ITaskService }      from "../api";


export class TaskController {
    private readonly taskService: ITaskService;

    constructor(taskService: ITaskService) {
        this.taskService = taskService;
    }

    async createTask(req: Request, res: Response): Promise<void> {
        res.json().send({ "status": "Ok" });
    }

    async getTaskById(req: Request, res: Response): Promise<void> {
        res.json().send({ "status": "Ok" });
    }

    async updateTask(req: Request, res: Response): Promise<void> {
        res.json().send({ "status": "Ok" });
    }

    async deleteTask(req: Request, res: Response): Promise<void> {
        res.json().send({ "status": "Ok" });
    }

    async getAllTasksForUser(req: Request, res: Response): Promise<void> {
        res.json().send({ "status": "Ok" });
    }

}