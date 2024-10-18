import { Request }                                from "express";
import * as jose                                  from "jose";
import { JWTPayload }                             from "jose";
import { UserCreationDto, UserDto, UserLoginDto } from "../dto";
import { User }                                   from "../model";


export interface IJwtService {
    extractClaims(token: string): Promise<JWTPayload>;

    isTokenExpired(claims: jose.JWTPayload): boolean;

    validateToken(token: string, username: string): Promise<boolean>;

    generateJwtToken(username: string): Promise<string>;
}

export interface IPasswordEncoder {
    matches(rawPassword: string, encodedPassword: string): boolean;

    encode(rawPassword: string): string;
}

export interface IAuthenticationService {
    loginUser(dto: UserLoginDto): Promise<UserDto>;

    registerUser(dto: UserCreationDto): Promise<UserDto>;
}

export interface AuthenticatedRequest extends Request {
    userInRequest?: User;
}