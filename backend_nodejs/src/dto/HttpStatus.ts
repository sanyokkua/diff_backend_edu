export enum HttpStatus {
    OK = 200,
    CREATED = 201,
    BAD_REQUEST = 400,
    UNAUTHORIZED = 401,
    FORBIDDEN = 403,
    NOT_FOUND = 404,
    CONFLICT = 409,
    NO_CONTENT = 204,
    INTERNAL_SERVER_ERROR = 500,
}

export namespace HttpStatus {
    export function getStatusMessage(statusCode: HttpStatus): string {
        switch (statusCode) {
            case HttpStatus.OK:
                return "OK";
            case HttpStatus.CREATED:
                return "Created";
            case HttpStatus.BAD_REQUEST:
                return "Bad Request";
            case HttpStatus.UNAUTHORIZED:
                return "Unauthorized";
            case HttpStatus.FORBIDDEN:
                return "Forbidden";
            case HttpStatus.NOT_FOUND:
                return "Not Found";
            case HttpStatus.CONFLICT:
                return "Conflict";
            case HttpStatus.NO_CONTENT:
                return "No Content";
            case HttpStatus.INTERNAL_SERVER_ERROR:
                return "Internal Server Error";
            default:
                return "Unknown Status";
        }
    }
}