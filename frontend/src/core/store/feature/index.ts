/**
 * Exports the default reducer from the GlobalSlice module.
 * @module GlobalReducer
 */
export { default as GlobalReducer } from "./global/GlobalSlice";

/**
 * Exports the default reducer from the UserSlice module.
 * @module UserReducer
 */
export { default as UserReducer }   from "./user/UserSlice";

/**
 * Exports the default reducer from the TaskSlice module.
 * @module TaskReducer
 */
export { default as TaskReducer }   from "./task/TaskSlice";

/**
 * Re-exports all named exports from the GlobalSlice module.
 * @module GlobalSlice
 */
export *                            from "./global/GlobalSlice";

/**
 * Re-exports all named exports from the UserSlice module.
 * @module UserSlice
 */
export *                            from "./user/UserSlice";

/**
 * Re-exports all named exports from the TaskSlice module.
 * @module TaskSlice
 */
export *                            from "./task/TaskSlice";
