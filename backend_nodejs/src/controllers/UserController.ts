import { Request, Response } from "express";
import { IUserService }      from "../api";


export class UserController {
    private readonly userService: IUserService;

    constructor(userService: IUserService) {
        this.userService = userService;
    }

    async getUserByID(req: Request, res: Response): Promise<void> {
        res.json().send({ "status": "Ok" });
    }

    async updateUserPassword(req: Request, res: Response): Promise<void> {
        res.json().send({ "status": "Ok" });
    }

    async deleteUser(req: Request, res: Response): Promise<void> {
        res.json().send({ "status": "Ok" });
    }

}