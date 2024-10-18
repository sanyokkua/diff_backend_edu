import { NextFunction, Response }                              from "express";
import { AuthenticatedRequest, IUserService }                  from "../api";
import { HttpStatus, UserDeletionDto, UserDto, UserUpdateDto } from "../dto";
import { ResponseDtoUtils }                                    from "../utils";
import { extractReqBody, extractUserFromContextAndValidate }   from "./Utils";


export class UserController {
    private readonly userService: IUserService;

    constructor(userService: IUserService) {
        this.userService = userService;

        // Bind methods
        this.getUserByID = this.getUserByID.bind(this);
        this.updateUserPassword = this.updateUserPassword.bind(this);
        this.deleteUser = this.deleteUser.bind(this);
    }

    async getUserByID(req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> {
        try {
            const userInRequest = extractUserFromContextAndValidate(req);
            const userDto: UserDto = {
                userId: userInRequest.id ?? -1,
                email: userInRequest.email,
                jwtToken: ""
            };
            ResponseDtoUtils.writePayloadToResponse(res, userDto, HttpStatus.OK);
        } catch (err) {
            next(err); // Pass error to the global error handler
        }
    }

    async updateUserPassword(req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> {
        try {
            const userInRequest = extractUserFromContextAndValidate(req);
            const userUpdateDto: UserUpdateDto = extractReqBody<UserUpdateDto>(req);
            const userDto: UserDto = await this.userService.updatePassword(userInRequest.id ?? -1, userUpdateDto);
            ResponseDtoUtils.writePayloadToResponse(res, userDto, HttpStatus.OK);
        } catch (err) {
            next(err); // Pass error to the global error handler
        }
    }

    async deleteUser(req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> {
        try {
            const userInRequest = extractUserFromContextAndValidate(req);
            const userDeletionDto: UserDeletionDto = extractReqBody<UserDeletionDto>(req);
            await this.userService.delete(userInRequest.id ?? -1, userDeletionDto);
            ResponseDtoUtils.writePayloadToResponse(res, "", HttpStatus.NO_CONTENT);
        } catch (err) {
            next(err); // Pass error to the global error handler
        }
    }
}
