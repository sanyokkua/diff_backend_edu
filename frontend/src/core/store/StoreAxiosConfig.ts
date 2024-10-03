import { AxiosClient, LogLevel } from "../config";
import { AppStore }              from "./ReduxStore";


const log = LogLevel.getLogger("StoreAxiosConfig");

/**
 * Configures Axios with the Redux store to include the JWT token in requests.
 * @function configureAxiosWithReduxStore
 * @param {AppStore} reduxStore - The Redux store instance.
 * @returns {void}
 */
export const configureAxiosWithReduxStore = (reduxStore: AppStore): void => {
    log.debug("configureAxiosWithReduxStore");

    AxiosClient.interceptors.request.use(
        async request => {
            const state = reduxStore.getState();
            const accessToken = state.users.userJwtToken;

            if (accessToken && request.headers) {
                request.headers.Authorization = `Bearer ${ accessToken }`;
            }

            return request;
        },
        async error => {
            return Promise.reject(error);
        }
    );

    AxiosClient.interceptors.response.use(
        response => {
            return response;
        },
        async error => {
            return Promise.reject(error);
        }
    );
};
