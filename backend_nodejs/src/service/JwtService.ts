import { createSecretKey }                from "crypto";
import { JWTPayload, jwtVerify, SignJWT } from "jose";
import { KeyObject }                      from "node:crypto";
import { IJwtService }                    from "../api";


const JWT_EXP_MINUTES = 15; // JWT expiration time in minutes

export class JwtService implements IJwtService {
    private readonly secretKey: KeyObject;

    constructor(secretKey: string) {
        this.secretKey = createSecretKey(Buffer.from(secretKey, "utf-8"));
    }

    async extractClaims(token: string): Promise<JWTPayload> {
        try {
            const { payload } = await jwtVerify(token, this.secretKey);
            return payload;
        } catch (err) {
            if ((err as any)?.name === "JWTExpired") {
                // Handle the expired token case separately if needed
                throw new Error("JWT token is expired");
            }
            throw new Error("Invalid JWT token");
        }
    }

    isTokenExpired(claims: JWTPayload): boolean {
        const exp = claims.exp;
        if (!exp) {
            console.warn("Token expiration date is missing.");
            return true; // Consider the token expired if the exp field is missing
        }

        const currentTime = Math.floor(Date.now() / 1000); // Current time in seconds
        return exp < currentTime;
    }

    async validateToken(token: string, username: string): Promise<boolean> {
        try {
            const claims = await this.extractClaims(token);
            const claimUsername = claims.sub;

            if (!claimUsername) {
                return false;
            }

            const isUsernameValid = claimUsername === username;
            const isTokenValid = !this.isTokenExpired(claims);

            return isUsernameValid && isTokenValid;
        } catch (err) {
            return false;
        }
    }

    async generateJwtToken(username: string): Promise<string> {
        const now = Math.floor(Date.now() / 1000);
        const exp = now + JWT_EXP_MINUTES * 60; // Expiration time in seconds

        return await new SignJWT({ sub: username })
            .setProtectedHeader({ alg: "HS256" })
            .setIssuedAt(now)
            .setExpirationTime(exp)
            .sign(this.secretKey);
    }

}