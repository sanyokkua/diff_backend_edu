import { AxiosClient, LogLevel } from "../config";
import { AppStore }              from "./ReduxStore";


const log = LogLevel.getLogger("StoreAxiosConfig");

export const configureAxiosWithReduxStore = (reduxStore: AppStore) => {
    log.debug("configureAxiosWithReduxStore");
    AxiosClient.interceptors.request.use(
        async request => {
            const state = reduxStore.getState();
            const accessToken = state.globals.userJwtToken;

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