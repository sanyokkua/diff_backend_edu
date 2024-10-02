import { StrictMode }              from "react";
import { createRoot }              from "react-dom/client";
import { Provider }                from "react-redux";
import { RouterProvider }          from "react-router-dom";
import { AppReduxStore, LogLevel } from "./core";

import { BrowserRouter } from "./react";


const log = LogLevel.getLogger("AppRoot");
const rootElement = document.getElementById("root");
if (!rootElement) {
    log.warn("Failed to load root element");
    throw new Error("Root element not found. Please ensure there is an element with id \"root\" in your HTML.");
}

log.debug("Will be rendered React App");
createRoot(rootElement).render(
    <StrictMode>
        <Provider store={ AppReduxStore }>
            <RouterProvider router={ BrowserRouter }/>
        </Provider>
    </StrictMode>
);