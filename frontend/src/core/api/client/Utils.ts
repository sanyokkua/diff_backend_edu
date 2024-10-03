import axios, { AxiosResponse } from "axios";
import { ResponseDto }          from "../dto/Common";


/**
 * Handles the response from an Axios request.
 * @template T
 * @param {AxiosResponse<ResponseDto<T>>} response - The Axios response object containing the data.
 * @returns {ResponseDto<T>} - The data transfer object containing the response data.
 * @throws Will throw an error if the response format is unexpected or data is missing.
 */
export const handleResponse = <T>(response: AxiosResponse<ResponseDto<T>>): ResponseDto<T> => {
    if (response?.data !== undefined) {
        return response.data;
    }
    throw new Error("Unexpected response format or missing data.");
};

/**
 * Handles errors from an Axios request.
 * @param {unknown} error - The error object thrown by Axios.
 * @throws Will throw an error with a message extracted from the Axios error response or a generic error message.
 */
export const handleError = (error: unknown): never => {
    if (axios.isAxiosError(error) && error.response) {
        if (error?.response?.data?.error) {
            throw new Error(error.response.data.error);
        } else {
            throw new Error(error.response.data.message || "An error occurred");
        }
    }
    throw new Error("An unexpected error occurred");
};

/**
 * Converts a timestamp in seconds to a Date object.
 * @param {number} seconds - The timestamp in seconds.
 * @returns {Date} - The corresponding Date object.
 */
export const getDateFromSeconds = (seconds: number): Date => {
    return new Date(seconds * 1000);
};

/**
 * Parses an error message from an unknown error object.
 * @param {unknown} error - The error object to parse.
 * @param {string} [defaultMsg="An unknown error occurred"] - The default message to return if the error object is not recognized.
 * @returns {string} - The parsed error message.
 */
export const parseErrorMessage = (error: unknown, defaultMsg: string = "An unknown error occurred"): string => {
    if (!error) {
        return defaultMsg;
    }

    if (typeof error === "string") {
        return error || defaultMsg;
    } else if (axios.isAxiosError(error)) {
        // Axios error
        if (error.response) {
            return error.response.data?.error || error.response.data?.message || error.message;
        } else if (error.request) {
            return "No response received from server";
        }
    } else if (error instanceof Error) {
        // General JavaScript/Node.js error
        return error.message;
    }

    return defaultMsg;
};
