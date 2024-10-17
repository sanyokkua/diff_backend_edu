import { Router }         from "express";
import { UserController } from "../controllers";


export const setupUserControllerRoutes = (userController: UserController): Router => {
    const router = Router();
    router.get("/:userId", userController.getUserByID);
    router.put("/:userId/password", userController.updateUserPassword);
    router.post("/:userId/delete", userController.deleteUser);

    return router;
};
