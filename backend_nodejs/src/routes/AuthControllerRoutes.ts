import { Router }         from "express";
import { AuthController } from "../controllers";


export const setupAuthControllerRoutes = (authController: AuthController): Router => {
    const router = Router();
    router.get("/login", authController.loginUser);
    router.put("/register", authController.registerUser);

    return router;
};
