import axios, { AxiosInstance } from "axios";


/**
 * Base URL for the API.
 * @constant {string}
 */
const baseUrl = "http://localhost:8080";

/**
 * Axios instance configured for the API.
 * This instance is pre-configured with the base URL, default headers, and credentials settings.
 * - `baseURL`: The base URL for the API endpoints.
 * - `headers`: Default headers for all requests, including `Content-Type` set to `application/json`.
 * - `withCredentials`: Indicates whether cross-site Access-Control requests should be made using credentials.
 * @type {axios.AxiosInstance}
 */
const axiosClient: AxiosInstance = axios.create(
    {
        baseURL: `${ baseUrl }/api/v1`,
        headers: {
            "Content-Type": "application/json"
        },
        withCredentials: true
    }
);

export default axiosClient;
