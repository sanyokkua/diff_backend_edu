import * as jose                                  from "jose";
import { UserCreationDTO, UserDTO, UserLoginDTO } from "../dto";


export interface IJwtService {
    extractClaims(token: string): jose.JWTPayload;

    isTokenExpired(claims: jose.JWTPayload): boolean;

    validateToken(token: string, username: string): boolean;

    generateJwtToken(username: string): string;
}

export interface IPasswordEncoder {
    matches(rawPassword: string, encodedPassword: string): boolean;

    encode(rawPassword: string): string;
}

export interface IAuthenticationService {
    loginUser(dto: UserLoginDTO): Promise<UserDTO>;

    registerUser(dto: UserCreationDTO): Promise<UserDTO>;
}

