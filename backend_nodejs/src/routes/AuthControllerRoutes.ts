import { Router }         from "express";
import { AuthController } from "../controller";


export const setupAuthControllerRoutes = (authController: AuthController): Router => {
    const router = Router({ mergeParams: true });
    router.post("/login", authController.loginUser.bind(authController));
    router.post("/register", authController.registerUser.bind(authController));

    return router;
};
