export type UserCreationDTO = {
    email: string;
    password: string;
    passwordConfirmation: string;
}

export type UserDeletionDTO = {
    email: string
    currentPassword: string
}

export type UserDTO = {
    userId: number
    email: string
    jwtToken: string
}

export type UserLoginDTO = {
    email: string
    password: string
}

export type UserUpdateDTO = {
    currentPassword: string
    newPassword: string
    newPasswordConfirmation: string
}
