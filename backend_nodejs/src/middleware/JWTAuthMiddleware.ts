import { NextFunction, Response }                                                  from "express";
import { AuthenticatedRequest, IJwtService, IUserRepository }                      from "../api";
import { HttpStatus }                                                              from "../dto";
import { AuthenticationCredentialsNotFoundError, InsufficientAuthenticationError } from "../error";
import { ResponseDtoUtils }                                                        from "../utils";


export class JWTAuthMiddleware {
    private readonly jwtService: IJwtService;
    private readonly userRepository: IUserRepository;

    constructor(jwtService: IJwtService, userRepository: IUserRepository) {
        this.jwtService = jwtService;
        this.userRepository = userRepository;
    }

    async process(req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> {
        console.log("req.params at the start of middleware:", req.params);
        const authHeader = req.header("Authorization");
        if (!authHeader) {
            ResponseDtoUtils.writeErrorToResponse(res, HttpStatus.UNAUTHORIZED, new InsufficientAuthenticationError("Auth Header absent"));
            return;
        }

        const tokenString = authHeader.replace("Bearer ", "").trim();
        if (!tokenString) {
            ResponseDtoUtils.writeErrorToResponse(res, HttpStatus.UNAUTHORIZED, new InsufficientAuthenticationError("Bearer token required"));
            return;
        }

        try {
            const claims = await this.jwtService.extractClaims(tokenString);
            const subject = claims.sub ?? "";

            const isValid: any = await this.jwtService.validateToken(tokenString, subject);
            if (!isValid) {
                ResponseDtoUtils.writeErrorToResponse(res, HttpStatus.UNAUTHORIZED, new InsufficientAuthenticationError("Invalid or expired token"));
                return;
            }

            const user = await this.userRepository.getUserByEmail(subject);
            if (!user) {
                ResponseDtoUtils.writeErrorToResponse(res, HttpStatus.UNAUTHORIZED, new AuthenticationCredentialsNotFoundError("User not found"));
                return;
            }

            req.userInRequest = user;
            next();
        } catch (err) {
            console.error("Error in JWT processing", err);
            ResponseDtoUtils.writeErrorToResponse(res, HttpStatus.UNAUTHORIZED, new InsufficientAuthenticationError("Invalid token"));
        }
    }
}
