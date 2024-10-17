import { Router }         from "express";
import { TaskController } from "../controllers";


export const setupTaskControllerRoutes = (taskController: TaskController): Router => {
    const router = Router();
    router.post("/", taskController.createTask);
    router.get("/:taskId", taskController.getTaskById);
    router.put("/:taskId", taskController.updateTask);
    router.delete("/:taskId", taskController.deleteTask);
    router.get("/", taskController.getAllTasksForUser);

    return router;
};
