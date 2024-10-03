import { StrictMode }              from "react";
import { createRoot }              from "react-dom/client";
import { Provider }                from "react-redux";
import { RouterProvider }          from "react-router-dom";
import { AppReduxStore, LogLevel } from "./core";
import { BrowserRouter }           from "./react";


const logger = LogLevel.getLogger("AppRoot");

/**
 * The root element of the application where the React app will be rendered.
 *
 * @type {HTMLElement | null}
 */
const rootElement: HTMLElement | null = document.getElementById("root");

if (!rootElement) {
    logger.warn("Failed to load root element");
    throw new Error("Root element not found. Please ensure there is an element with id \"root\" in your HTML.");
}

logger.debug("Will be rendered React App");

/**
 * Renders the React application into the root element.
 *
 * @function
 */
createRoot(rootElement).render(
    <StrictMode>
        <Provider store={ AppReduxStore }>
            <RouterProvider router={ BrowserRouter }/>
        </Provider>
    </StrictMode>
);
