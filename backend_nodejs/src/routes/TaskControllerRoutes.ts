import { Router }            from "express";
import { TaskController }    from "../controller";
import { JWTAuthMiddleware } from "../middleware/JWTAuthMiddleware";


export const setupTaskControllerRoutes = (taskController: TaskController, jwtMiddleware: JWTAuthMiddleware): Router => {
    const router = Router({ mergeParams: true });
    router.use(jwtMiddleware.process.bind(jwtMiddleware));
    router.post("/", taskController.createTask.bind(taskController));
    router.get("/:taskId", taskController.getTaskById.bind(taskController));
    router.put("/:taskId", taskController.updateTask.bind(taskController));
    router.delete("/:taskId", taskController.deleteTask.bind(taskController));
    router.get("/", taskController.getAllTasksForUser.bind(taskController));

    return router;
};
