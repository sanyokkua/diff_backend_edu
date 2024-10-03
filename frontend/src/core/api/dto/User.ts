/**
 * Data Transfer Object for user data.
 *
 * @property {number} userId - The unique identifier of the user.
 * @property {string} email - The email of the user.
 * @property {string} [jwtToken] - The JSON Web Token (JWT) for the user, if available.
 */
export interface UserDto {
    userId: number;
    email: string;
    jwtToken?: string;
}

/**
 * Data Transfer Object for user deletion.
 *
 * @property {string} email - The email of the user to be deleted.
 * @property {string} currentPassword - The current password of the user for verification.
 */
export interface UserDeletionDTO {
    email: string;
    currentPassword: string;
}

/**
 * Data Transfer Object for user update.
 *
 * @property {string} currentPassword - The current password of the user for verification.
 * @property {string} newPassword - The new password for the user.
 * @property {string} newPasswordConfirmation - The confirmation of the new password.
 */
export interface UserUpdateDTO {
    currentPassword: string;
    newPassword: string;
    newPasswordConfirmation: string;
}
