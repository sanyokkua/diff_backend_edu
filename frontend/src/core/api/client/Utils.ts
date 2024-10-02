import axios, { AxiosResponse } from "axios";
import { ResponseDto }          from "../dto/Common";


export const handleResponse = <T>(response: AxiosResponse<ResponseDto<T>>): ResponseDto<T> => {
    if (response?.data !== undefined) {
        return response.data;
    }
    throw new Error("Unexpected response format or missing data.");
};

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

export const getDateFromSeconds = (seconds: number): Date => {
    return new Date(seconds * 1000);
};

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
