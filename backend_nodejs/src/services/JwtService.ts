import { JWTPayload }  from "jose";
import { IJwtService } from "../api";


export class JwtService implements IJwtService {
    extractClaims(token: string): JWTPayload {
        throw new Error("Method not implemented.");
    }

    isTokenExpired(claims: JWTPayload): boolean {
        throw new Error("Method not implemented.");
    }

    validateToken(token: string, username: string): boolean {
        throw new Error("Method not implemented.");
    }

    generateJwtToken(username: string): string {
        throw new Error("Method not implemented.");
    }

}