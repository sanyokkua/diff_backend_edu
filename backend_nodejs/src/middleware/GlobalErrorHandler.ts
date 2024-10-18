import { NextFunction, Request, Response } from "express";
import { HttpStatus }                      from "../dto";

import {
    AccessDeniedError,
    AuthenticationCredentialsNotFoundError,
    EmailAlreadyExistsError,
    IllegalArgumentError,
    InsufficientAuthenticationError,
    InvalidEmailFormatError,
    InvalidJwtTokenError,
    InvalidPasswordError,
    NoHandlerFoundError,
    TaskAlreadyExistsError,
    TaskNotFoundError
}                           from "../error";
import { ResponseDtoUtils } from "../utils";


function getErrorStatusCode(err: Error, statusDefault: HttpStatus = HttpStatus.INTERNAL_SERVER_ERROR): HttpStatus {
    console.debug("Entering getErrorStatusCode");

    if (err instanceof EmailAlreadyExistsError) {
        return HttpStatus.CONFLICT;
    } else if (err instanceof InvalidEmailFormatError) {
        return HttpStatus.BAD_REQUEST;
    } else if (err instanceof InvalidJwtTokenError) {
        return HttpStatus.BAD_REQUEST;
    } else if (err instanceof InvalidPasswordError) {
        return HttpStatus.BAD_REQUEST;
    } else if (err instanceof TaskAlreadyExistsError) {
        return HttpStatus.CONFLICT;
    } else if (err instanceof TaskNotFoundError) {
        return HttpStatus.NOT_FOUND;
    } else if (err instanceof AccessDeniedError) {
        return HttpStatus.FORBIDDEN;
    } else if (err instanceof AuthenticationCredentialsNotFoundError) {
        return HttpStatus.UNAUTHORIZED;
    } else if (err instanceof InsufficientAuthenticationError) {
        return HttpStatus.UNAUTHORIZED;
    } else if (err instanceof NoHandlerFoundError) {
        return HttpStatus.NOT_FOUND;
    } else if (err instanceof IllegalArgumentError) {
        return HttpStatus.BAD_REQUEST;
    } else {
        return statusDefault;
    }
}

export const globalErrorHandler = (err: Error, req: Request, res: Response, next: NextFunction) => {
    console.error(err.stack);
    const statusCode = getErrorStatusCode(err);
    ResponseDtoUtils.writeErrorToResponse(res, statusCode, err);
};