/**
 * Data Transfer Object for user login.
 *
 * @property {string} email - The email of the user.
 * @property {string} password - The password of the user.
 */
export interface UserLoginDto {
    email: string;
    password: string;
}

/**
 * Data Transfer Object for user creation.
 *
 * @property {string} email - The email of the user.
 * @property {string} password - The password of the user.
 * @property {string} passwordConfirmation - The confirmation of the user's password.
 */
export interface UserCreationDTO {
    email: string;
    password: string;
    passwordConfirmation: string;
}
