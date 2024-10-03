import { TypedUseSelectorHook, useDispatch, useSelector } from "react-redux";
import type { AppDispatch, RootState }                    from "./ReduxStore";


/**
 * Custom hook to dispatch actions in the Redux store.
 * @function useAppDispatch
 * @returns {AppDispatch} - The dispatch function for the Redux store.
 */
export const useAppDispatch = (): AppDispatch => useDispatch<AppDispatch>();

/**
 * Custom hook to select state from the Redux store.
 * @type {TypedUseSelectorHook<RootState>}
 */
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;
