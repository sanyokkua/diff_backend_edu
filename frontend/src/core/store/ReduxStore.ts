import { configureStore }                          from "@reduxjs/toolkit";
import { LogLevel }                                from "../config";
import BrowserStore, { ItemType }                  from "./BrowserStore";
import { GlobalReducer, TaskReducer, UserReducer } from "./feature";
import { configureAxiosWithReduxStore }            from "./StoreAxiosConfig";


const log = LogLevel.getLogger("ReduxStore");

/**
 * Configures the Redux store with the specified reducers.
 * @constant AppReduxStore
 */
const AppReduxStore = configureStore(
    {
        reducer: {
            globals: GlobalReducer,
            users: UserReducer,
            tasks: TaskReducer
        }
    }
);

let previousState = AppReduxStore.getState();

/**
 * Subscribes to store updates and saves relevant parts of the state to local storage.
 * @function
 */
AppReduxStore.subscribe(() => {
    try {
        const currentState = AppReduxStore.getState();
        const storeApi = new BrowserStore();

        const prevUserId = previousState.users.userId;
        const prevUserEmail = previousState.users.userEmail;
        const prevUserJwt = previousState.users.userJwtToken;

        const currentUserId = currentState.users.userId;
        const currentUserEmail = currentState.users.userEmail;
        const currentUserJwt = currentState.users.userJwtToken;

        if (prevUserId !== currentUserId) {
            storeApi.saveData(ItemType.USER_ID, `${ currentUserId }`);
        }
        if (prevUserEmail !== currentUserEmail) {
            storeApi.saveData(ItemType.USER_EMAIL, currentUserEmail);
        }
        if (prevUserJwt !== currentUserJwt) {
            storeApi.saveData(ItemType.JWT_TOKEN, currentUserJwt);
        }

        previousState = currentState;
    } catch (err) {
        log.error("Could not save partial state", err);
    }
});

/**
 * Type representing the root state of the Redux store.
 */
export type RootState = ReturnType<typeof AppReduxStore.getState>;

/**
 * Type representing the dispatch function of the Redux store.
 */
export type AppDispatch = typeof AppReduxStore.dispatch;

/**
 * Type representing the Redux store.
 */
export type AppStore = typeof AppReduxStore;

configureAxiosWithReduxStore(AppReduxStore);

export default AppReduxStore;
