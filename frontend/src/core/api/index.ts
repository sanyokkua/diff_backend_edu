/**
 * Exports the default AuthClient class from the AuthClient module.
 * @module AuthClient
 */
export { default as AuthClient }                                              from "./client/AuthClient";

/**
 * Exports the default TaskClient class from the TaskClient module.
 * @module TaskClient
 */
export { default as TaskClient }                                              from "./client/TaskClient";

/**
 * Exports the default UserClient class from the UserClient module.
 * @module UserClient
 */
export { default as UserClient }                                              from "./client/UserClient";

/**
 * Exports utility functions from the Utils module.
 * @module Utils
 */
export { handleError, handleResponse, parseErrorMessage, getDateFromSeconds } from "./client/Utils";

/**
 * Exports everything from the LoginSchema module.
 * @module LoginSchema
 */
export *                                                                      from "./schema/LoginSchema";

/**
 * Exports everything from the RegistrationSchema module.
 * @module RegistrationSchema
 */
export *                                                                      from "./schema/RegistrationSchema";

/**
 * Exports everything from the SchemaUtils module.
 * @module SchemaUtils
 */
export *                                                                      from "./schema/SchemaUtils";

/**
 * Exports types related to user authentication from the Auth module.
 * @module Auth
 */
export type { UserCreationDTO, UserLoginDto }                                 from "./dto/Auth";

/**
 * Exports the ResponseDto type from the Common module.
 * @module Common
 */
export type { ResponseDto }                                                   from "./dto/Common";

/**
 * Exports types related to tasks from the Task module.
 * @module Task
 */
export type { TaskCreationDTO, TaskDto, TaskUpdateDTO }                       from "./dto/Task";

/**
 * Exports types related to user management from the User module.
 * @module User
 */
export type { UserDeletionDTO, UserDto, UserUpdateDTO }                       from "./dto/User";
