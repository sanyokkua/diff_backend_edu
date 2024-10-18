import { NextFunction, Response }                                           from "express";
import { AuthenticatedRequest, ITaskService }                               from "../api";
import { HttpStatus, TaskCreationDto, TaskDto, TaskUpdateDto }              from "../dto";
import { ResponseDtoUtils }                                                 from "../utils";
import { extractReqBody, extractTaskId, extractUserFromContextAndValidate } from "./Utils";


export class TaskController {
    private readonly taskService: ITaskService;

    constructor(taskService: ITaskService) {
        this.taskService = taskService;

        // Bind methods
        this.createTask = this.createTask.bind(this);
        this.getTaskById = this.getTaskById.bind(this);
        this.updateTask = this.updateTask.bind(this);
        this.deleteTask = this.deleteTask.bind(this);
        this.getAllTasksForUser = this.getAllTasksForUser.bind(this);
    }

    async createTask(req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> {
        try {
            console.log("req.params:", req.params);  // Debug line
            const userInRequest = extractUserFromContextAndValidate(req);
            const taskCreationDto: TaskCreationDto = extractReqBody<TaskCreationDto>(req);
            const taskDto: TaskDto = await this.taskService.createTask(userInRequest.id ?? -1, taskCreationDto);
            ResponseDtoUtils.writePayloadToResponse(res, taskDto, HttpStatus.CREATED);
        } catch (err) {
            next(err); // Pass error to the global error handler
        }
    }

    async getTaskById(req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> {
        try {
            console.log("req.params:", req.params);  // Debug line
            const userInRequest = extractUserFromContextAndValidate(req);
            const taskId = extractTaskId(req);
            const taskDto: TaskDto = await this.taskService.getTaskByUserIDAndTaskID(userInRequest.id ?? -1, taskId);
            ResponseDtoUtils.writePayloadToResponse(res, taskDto, HttpStatus.OK);
        } catch (err) {
            next(err); // Pass error to the global error handler
        }
    }

    async updateTask(req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> {
        try {
            const userInRequest = extractUserFromContextAndValidate(req);
            const taskId = extractTaskId(req);
            const taskUpdateDto: TaskUpdateDto = extractReqBody<TaskUpdateDto>(req);
            const taskDto: TaskDto = await this.taskService.updateTask(userInRequest.id ?? -1, taskId, taskUpdateDto);
            ResponseDtoUtils.writePayloadToResponse(res, taskDto, HttpStatus.OK);
        } catch (err) {
            next(err); // Pass error to the global error handler
        }
    }

    async deleteTask(req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> {
        try {
            const userInRequest = extractUserFromContextAndValidate(req);
            const taskId = extractTaskId(req);
            await this.taskService.deleteTask(userInRequest.id ?? -1, taskId);
            ResponseDtoUtils.writePayloadToResponse(res, "", HttpStatus.NO_CONTENT);
        } catch (err) {
            next(err); // Pass error to the global error handler
        }
    }

    async getAllTasksForUser(req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> {
        try {
            console.log("req.params:", req.params);  // Debug line
            const userInRequest = extractUserFromContextAndValidate(req);
            const taskDtos: TaskDto[] = await this.taskService.getAllTasksForUser(userInRequest.id ?? -1);
            ResponseDtoUtils.writePayloadToResponse(res, taskDtos, HttpStatus.OK);
        } catch (err) {
            next(err); // Pass error to the global error handler
        }
    }
}
