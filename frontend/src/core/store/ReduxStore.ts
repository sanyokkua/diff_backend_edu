import { configureStore }               from "@reduxjs/toolkit";
import { LogLevel }                     from "../config";
import BrowserStore, { ItemType }       from "./BrowserStore";
import { GlobalSlice }                  from "./feature";
import { configureAxiosWithReduxStore } from "./StoreAxiosConfig";


const log = LogLevel.getLogger("ReduxStore");

const AppReduxStore = configureStore({ reducer: { globals: GlobalSlice } });

AppReduxStore.subscribe(() => {
    try {
        const storeApi = new BrowserStore();
        const state = AppReduxStore.getState();

        const userId = state.globals.userId;
        const userEmail = state.globals.userEmail;
        const userJwt = state.globals.userJwtToken;

        storeApi.saveData(ItemType.USER_ID, `${ userId }`);
        storeApi.saveData(ItemType.USER_EMAIL, userEmail);
        storeApi.saveData(ItemType.JWT_TOKEN, userJwt);
    } catch (err) {
        log.error("Could not save partial state", err);
    }
});

export type RootState = ReturnType<typeof AppReduxStore.getState>;
export type AppDispatch = typeof AppReduxStore.dispatch;
export type AppStore = typeof AppReduxStore;

configureAxiosWithReduxStore(AppReduxStore);

export default AppReduxStore;