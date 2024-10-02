export interface UserLoginDto {
    email: string;
    password: string;
}

export interface UserCreationDTO {
    email: string;
    password: string;
    passwordConfirmation: string;
}