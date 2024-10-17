import bodyParser                                                                               from "body-parser";
import express, { Application }                                                                 from "express";
import {
    getDataSource
}                                                                                               from "./config/Database";
import { AuthController, TaskController, UserController }                                       from "./controllers";
import {
    globalErrorHandler
}                                                                                               from "./middleware/GlobalErrorHandler";
import { Task, User }                                                                           from "./models";
import { TaskRepository, UserRepository }                                                       from "./repository";
import {
    setupAuthControllerRoutes
}                                                                                               from "./routes/AuthControllerRoutes";
import {
    setupTaskControllerRoutes
}                                                                                               from "./routes/TaskControllerRoutes";
import {
    setupUserControllerRoutes
}                                                                                               from "./routes/UserControllerRoutes";
import { AuthenticationService, BCryptPasswordEncoder, DEFAULT_COST, TaskService, UserService } from "./services";
import {
    JwtService
}                                                                                               from "./services/JwtService";


export const createExpressApp = (): Application => {
    const datasource = getDataSource();
    const userOrmRepository = datasource.getRepository(User);
    const taskOrmRepository = datasource.getRepository(Task);

    const userRepo = new UserRepository(userOrmRepository);
    const taskRepo = new TaskRepository(taskOrmRepository);
    const passwordEncoder = new BCryptPasswordEncoder(DEFAULT_COST);
    const userService = new UserService(userRepo, passwordEncoder);
    const taskService = new TaskService(taskRepo, userRepo);
    const jwtService = new JwtService();
    const authService = new AuthenticationService(userService, userRepo, jwtService, passwordEncoder);
    const authController = new AuthController(authService);
    const userController = new UserController(userService);
    const taskController = new TaskController(taskService);

    const authRouter = setupAuthControllerRoutes(authController);
    const userRouter = setupUserControllerRoutes(userController);
    const taskRouter = setupTaskControllerRoutes(taskController);

    const app: Application = express();
    app.use(bodyParser.json());
    app.use(globalErrorHandler);
    app.use("/api/v1/auth", authRouter);
    app.use("/api/v1/users", userRouter);
    app.use("/api/v1/users/:userId/tasks", taskRouter);

    return app;
};
