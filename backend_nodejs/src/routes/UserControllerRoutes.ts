import { Router }            from "express";
import { UserController }    from "../controller";
import { JWTAuthMiddleware } from "../middleware/JWTAuthMiddleware";


export const setupUserControllerRoutes = (userController: UserController, jwtMiddleware: JWTAuthMiddleware): Router => {
    const router = Router({ mergeParams: true });
    router.use(jwtMiddleware.process.bind(jwtMiddleware));
    router.get("/:userId", userController.getUserByID.bind(userController));
    router.put("/:userId/password", userController.updateUserPassword.bind(userController));
    router.post("/:userId/delete", userController.deleteUser.bind(userController));

    return router;
};
