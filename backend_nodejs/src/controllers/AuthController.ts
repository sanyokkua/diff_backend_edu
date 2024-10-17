import { Request, Response }      from "express";
import { IAuthenticationService } from "../api";


export class AuthController {
    private readonly authenticationService: IAuthenticationService;

    constructor(authenticationService: IAuthenticationService) {
        this.authenticationService = authenticationService;
    }

    async loginUser(req: Request, res: Response): Promise<void> {
        res.json().send({ "status": "Ok" });
    }

    async registerUser(req: Request, res: Response): Promise<void> {
        res.json().send({ "status": "Ok" });
    }

}