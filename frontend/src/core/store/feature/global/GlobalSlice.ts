import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { FeedbackType }               from "../../../../react";
import { LogLevel }                   from "../../../config";


const logger = LogLevel.getLogger("GlobalSlice");

/**
 * Represents the global state of the application.
 *
 * @property {string} headerTitle - The title displayed in the app header.
 * @property {boolean} isLoading - Indicates whether the app is currently loading.
 * @property {FeedbackType} feedback - The feedback message and its severity.
 */
export interface GlobalState {
    headerTitle: string;
    isLoading: boolean;
    feedback: FeedbackType;
}

const initialState: GlobalState = {
    headerTitle: "",
    isLoading: false,
    feedback: {
        message: "",
        severity: "success"
    }
};

/**
 * Slice for managing global state in the application.
 * @namespace globalSlice
 */
export const globalSlice = createSlice(
    {
        name: "global",
        initialState,
        reducers: {
            /**
             * Sets the header title of the application.
             * @function setHeaderTitle
             * @memberof globalSlice
             * @param state - The current state of the global slice.
             * @param action - The action containing the new header title.
             */
            setHeaderTitle: (state, action: PayloadAction<string>) => {
                const newTitle = action.payload;
                logger.debug(`Updating app header title to: "${ newTitle }"`);
                state.headerTitle = newTitle;
            },

            /**
             * Sets the loading state of the application.
             * @function setIsLoading
             * @memberof globalSlice
             * @param state - The current state of the global slice.
             * @param {PayloadAction<boolean>} action - The action containing the new loading state.
             */
            setIsLoading: (state, action: PayloadAction<boolean>) => {
                const isLoading = action.payload;
                logger.debug(`Setting app loading state to: ${ isLoading }`);
                state.isLoading = isLoading;
            },

            /**
             * Sets the feedback message and severity.
             * @function setFeedback
             * @memberof globalSlice
             * @param state - The current state of the global slice.
             * @param action - The action containing the new feedback.
             */
            setFeedback: (state, action: PayloadAction<FeedbackType>) => {
                const feedback = action.payload;
                logger.error(`Updating feedback to: "${ feedback }"`);
                state.feedback = feedback;
            },

            /**
             * Clears the feedback message.
             * @function clearFeedback
             * @memberof globalSlice
             * @param state - The current state of the global slice.
             */
            clearFeedback: (state) => {
                logger.debug("Clearing app error message");
                state.feedback = {
                    message: "",
                    severity: "success"
                };
            }
        }
    });

/**
 * User slice action creators.
 *
 * @property {function} setHeaderTitle - Action to set the HeaderTitle.
 * @property {function} setIsLoading - Action to set the IsLoading.
 * @property {function} setFeedback - Action to set the Feedback.
 * @property {function} clearFeedback - Action to clear the Feedback.
 */
export const {
    setHeaderTitle,
    setIsLoading,
    setFeedback,
    clearFeedback
} = globalSlice.actions;

export default globalSlice.reducer;
