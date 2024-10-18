import { Response }                from "express";
import { HttpStatus, ResponseDto } from "../dto";


export class ResponseDtoUtils {

    static getErrorMessage(ex: Error): string {
        const type = ex.name;
        const message = ex.message;
        return `${ type }: ${ message }`;
    }

    static buildDtoResponse<T>(data: T, status: HttpStatus): ResponseDto<T> {
        const statusCode = status;
        const statusMessage = HttpStatus.getStatusMessage(status);
        return {
            data,
            statusCode,
            statusMessage,
            error: ""
        };
    }

    static writePayloadToResponse<T>(res: Response, data: T, status: HttpStatus): void {
        const buildDtoResponse: ResponseDto<T> = ResponseDtoUtils.buildDtoResponse(data, status);
        res.status(status).json(buildDtoResponse);
    }

    static buildDtoErrorResponse<T>(status: HttpStatus, ex: Error): ResponseDto<T> {
        return ResponseDtoUtils.createErrorResponseBody(status, ex);
    }

    static writeErrorToResponse<T>(res: Response, status: HttpStatus, ex: Error): void {
        const buildDtoResponse: ResponseDto<T> = ResponseDtoUtils.buildDtoErrorResponse(status, ex);
        res.status(status).json(buildDtoResponse);
    }

    static createErrorResponseBody<T>(status: HttpStatus, ex: Error): ResponseDto<T> {
        const msg = ResponseDtoUtils.getErrorMessage(ex);
        const statusCode = status;
        const statusMessage = HttpStatus.getStatusMessage(status);
        return {
            statusCode,
            statusMessage,
            error: msg
        };
    }
}
