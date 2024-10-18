import { IPasswordEncoder, IUserRepository }                             from "../api";
import { UserCreationDto, UserDeletionDto, UserLoginDto, UserUpdateDto } from "../dto";
import {
    AccessDeniedError,
    EmailAlreadyExistsError,
    IllegalArgumentError,
    InvalidEmailFormatError,
    InvalidPasswordError
}                                                                        from "../error";
import { User }                                                          from "../model";

// Email validation regex
const EmailPattern = /^[\w-.]+@([\w-]+\.)[\w-]{2,4}$/;

export function validateEmailFormat(email: string): void {
    if (email === "" || !EmailPattern.test(email)) {
        throw new InvalidEmailFormatError("Invalid email format");
    }
}

export async function checkUserExists(userRepo: IUserRepository, email: string): Promise<void> {
    try {
        const user = await userRepo.getUserByEmail(email);
        if (user) {
            throw new EmailAlreadyExistsError("Email already exists");
        }
    } catch (err: unknown) {
        if ((err as any)?.message !== "RecordNotFound") { //TODO: verify
            throw err;
        }
    }
}

export function validatePasswords(password: string, passwordConfirmation: string): void {
    if (password === "" || passwordConfirmation === "") {
        throw new InvalidPasswordError("Passwords can't have empty value");
    }
    if (password !== passwordConfirmation) {
        throw new InvalidPasswordError("Passwords do not match");
    }
}

export function validatePasswordUpdate(userUpdateDTO: UserUpdateDto, user: User, passwordEncoder: IPasswordEncoder): void {
    const matches = passwordEncoder.matches(userUpdateDTO.currentPassword, user.passwordHash);
    if (!matches) {
        throw new InvalidPasswordError("Current password is incorrect");
    }

    if (userUpdateDTO.newPassword === userUpdateDTO.currentPassword) {
        throw new InvalidPasswordError("New password cannot be the same as the current password");
    }

    validatePasswords(userUpdateDTO.newPassword, userUpdateDTO.newPasswordConfirmation);
}

export function validateAuthenticatedUserID(userIdFromDto: number, userID: number): void {
    if (userIdFromDto === 0 || userID === 0 || userIdFromDto !== userID) {
        throw new AccessDeniedError("User is not authorized to perform this action");
    }
}

export function validateUserLoginDto(userLoginDto: UserLoginDto | null): void {
    if (userLoginDto === null) {
        throw new IllegalArgumentError("UserLoginDto is nil");
    }
    if (userLoginDto.email === "") {
        throw new IllegalArgumentError("UserLoginDto email is nil or empty");
    }
    if (userLoginDto.password === "") {
        throw new IllegalArgumentError("UserLoginDto password is nil or empty");
    }
}

export function validateUserCreationDTO(userCreationDTO: UserCreationDto | null): void {
    if (userCreationDTO === null) {
        throw new IllegalArgumentError("UserCreationDto is nil");
    }
    if (userCreationDTO.email === "") {
        throw new IllegalArgumentError("UserCreationDto email is nil or empty");
    }
    if (userCreationDTO.password === "") {
        throw new IllegalArgumentError("UserCreationDto password is nil or empty");
    }
    if (userCreationDTO.passwordConfirmation === "") {
        throw new IllegalArgumentError("UserCreationDto password confirmation is nil or empty");
    }
    validatePasswords(userCreationDTO.password, userCreationDTO.passwordConfirmation);
}

export function validateUserUpdateDTO(userUpdateDTO: UserUpdateDto | null): void {
    if (userUpdateDTO === null) {
        throw new IllegalArgumentError("UserUpdateDto is nil");
    }
    if (userUpdateDTO.currentPassword === "") {
        throw new IllegalArgumentError("UserUpdateDto currentPassword is nil or empty");
    }
    if (userUpdateDTO.newPassword === "") {
        throw new IllegalArgumentError("UserUpdateDto newPassword is nil or empty");
    }
    if (userUpdateDTO.newPasswordConfirmation === "") {
        throw new IllegalArgumentError("UserUpdateDto newPasswordConfirmation is nil or empty");
    }
    if (userUpdateDTO.newPasswordConfirmation !== userUpdateDTO.newPassword) {
        throw new IllegalArgumentError("Passwords do not match");
    }
}

export function validateUserDeletionDTO(userDeletionDTO: UserDeletionDto | null): void {
    if (userDeletionDTO === null) {
        throw new IllegalArgumentError("UserDeletionDto is nil");
    }
    if (userDeletionDTO.email === "") {
        throw new IllegalArgumentError("UserDeletionDto email is nil or empty");
    }
    if (userDeletionDTO.currentPassword === "") {
        throw new IllegalArgumentError("UserDeletionDto password is nil or empty");
    }
}
