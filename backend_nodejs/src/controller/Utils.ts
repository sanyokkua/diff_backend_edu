import { Request }                                                      from "express";
import { AuthenticatedRequest }                                         from "../api";
import { AuthenticationCredentialsNotFoundError, IllegalArgumentError } from "../error";
import { User }                                                         from "../model";
import { validateAuthenticatedUserID }                                  from "../utils";


export function extractUserId(req: Request): number {
    console.log(req.params);
    const id = (req.params as any)?.userId;
    if (!id) {
        throw new IllegalArgumentError("User ID is not found in path");
    }
    return Number.parseInt(id);
}

export function extractTaskId(req: Request): number {
    console.log(req.params);
    const id = (req.params as any)?.taskId;
    if (!id) {
        throw new IllegalArgumentError("Task ID is not found in path");
    }
    return Number.parseInt(id);
}

export function extractReqBody<T>(req: Request): T {
    console.log(req.body);
    const body: T = req.body as T;
    if (!body) {
        throw new IllegalArgumentError("Request doesn't have request body");
    }
    return body;
}

export function extractUserFromContextAndValidate(req: AuthenticatedRequest): User {
    console.log(req.params);
    const userId: number = extractUserId(req);
    const userInRequest: User | undefined = req.userInRequest;
    if (!userInRequest) {
        throw new AuthenticationCredentialsNotFoundError("User not logged in");
    }
    validateAuthenticatedUserID(userInRequest.id ?? -1, userId);

    return userInRequest;
}
