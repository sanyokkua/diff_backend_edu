export interface ResponseDto<T> {
    statusCode: null;
    statusMessage: string;
    data: T;
    error: string;
}
