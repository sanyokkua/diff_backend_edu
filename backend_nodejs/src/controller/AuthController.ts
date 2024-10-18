import { NextFunction, Request, Response }                    from "express";
import { IAuthenticationService }                             from "../api";
import { HttpStatus, UserCreationDto, UserDto, UserLoginDto } from "../dto";
import { ResponseDtoUtils }                                   from "../utils";
import { extractReqBody }                                     from "./Utils";


export class AuthController {
    private readonly authenticationService: IAuthenticationService;

    constructor(authenticationService: IAuthenticationService) {
        this.authenticationService = authenticationService;

        // Bind methods
        this.loginUser = this.loginUser.bind(this);
        this.registerUser = this.registerUser.bind(this);
    }

    async loginUser(req: Request, res: Response, next: NextFunction): Promise<void> {
        try {
            const userLoginDto: UserLoginDto = extractReqBody<UserLoginDto>(req);
            const userDto: UserDto = await this.authenticationService.loginUser(userLoginDto);
            ResponseDtoUtils.writePayloadToResponse(res, userDto, HttpStatus.OK);
        } catch (err) {
            next(err); // Pass error to the global error handler
        }
    }

    async registerUser(req: Request, res: Response, next: NextFunction): Promise<void> {
        try {
            const userRegistrationDto: UserCreationDto = extractReqBody<UserCreationDto>(req);
            const userDto: UserDto = await this.authenticationService.registerUser(userRegistrationDto);
            ResponseDtoUtils.writePayloadToResponse(res, userDto, HttpStatus.CREATED);
        } catch (err) {
            next(err); // Pass error to the global error handler
        }
    }
}
