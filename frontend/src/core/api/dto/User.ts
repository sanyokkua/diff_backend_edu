export interface UserDto {
    userId: number;
    email: string;
    jwtToken?: string;
}

export interface UserDeletionDTO {
    email: string;
    currentPassword: string;
}

export interface UserUpdateDTO {
    currentPassword: string;
    newPassword: string;
    newPasswordConfirmation: string;
}

