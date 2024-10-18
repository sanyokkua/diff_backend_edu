export type ResponseDto<T> = {
    statusCode: number;
    statusMessage: string;
    data?: T;
    error: string;
}