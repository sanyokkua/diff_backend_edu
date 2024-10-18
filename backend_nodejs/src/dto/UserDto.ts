export type UserCreationDto = {
    email: string;
    password: string;
    passwordConfirmation: string;
}

export type UserDeletionDto = {
    email: string
    currentPassword: string
}

export type UserDto = {
    userId: number
    email: string
    jwtToken: string
}

export type UserLoginDto = {
    email: string
    password: string
}

export type UserUpdateDto = {
    currentPassword: string
    newPassword: string
    newPasswordConfirmation: string
}
