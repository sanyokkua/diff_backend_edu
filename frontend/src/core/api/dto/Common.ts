/**
 * Generic Data Transfer Object for API responses.
 *
 * @property {null} statusCode - The status code of the response. Currently set to null.
 * @property {string} statusMessage - The status message of the response.
 * @property data - The data returned by the API.
 * @property {string} error - The error message, if any.
 */
export interface ResponseDto<T> {
    statusCode: null;
    statusMessage: string;
    data: T;
    error: string;
}
