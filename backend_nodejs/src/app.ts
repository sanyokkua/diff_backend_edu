import bodyParser                                                                               from "body-parser";
import cors                                                                                     from "cors";
import express, { Application }                                                                 from "express";
import {
    getDataSource
}                                                                                               from "./config/Database";
import { AuthController, TaskController, UserController }                                       from "./controller";
import {
    globalErrorHandler
}                                                                                               from "./middleware/GlobalErrorHandler";
import {
    JWTAuthMiddleware
}                                                                                               from "./middleware/JWTAuthMiddleware";
import { Task, User }                                                                           from "./model";
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
import { AuthenticationService, BCryptPasswordEncoder, DEFAULT_COST, TaskService, UserService } from "./service";
import {
    JwtService
}                                                                                               from "./service/JwtService";


const corsOptions = {
    origin: "http://localhost:5173",
    methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
    allowedHeaders: ["Origin", "Content-Type", "Authorization"],
    exposedHeaders: ["Content-Length"],
    credentials: true,
    maxAge: 12 * 60 * 60 // 12 hours
};

export const createExpressApp = async (): Promise<Application> => {
    const datasource = getDataSource();
    await datasource.initialize();

    const userOrmRepository = datasource.getRepository(User);
    const taskOrmRepository = datasource.getRepository(Task);

    const userRepo = new UserRepository(userOrmRepository);
    const taskRepo = new TaskRepository(taskOrmRepository);
    const passwordEncoder = new BCryptPasswordEncoder(DEFAULT_COST);
    const userService = new UserService(userRepo, passwordEncoder);
    const taskService = new TaskService(taskRepo, userRepo);
    const jwtService = new JwtService(process.env.JWT_SECRET ?? "");
    const authService = new AuthenticationService(userService, userRepo, jwtService, passwordEncoder);
    const authController = new AuthController(authService);
    const userController = new UserController(userService);
    const taskController = new TaskController(taskService);
    const jwtAuthMiddleware = new JWTAuthMiddleware(jwtService, userRepo);

    const authRouter = setupAuthControllerRoutes(authController);
    const userRouter = setupUserControllerRoutes(userController, jwtAuthMiddleware);
    const taskRouter = setupTaskControllerRoutes(taskController, jwtAuthMiddleware);

    const app: Application = express();
    app.use(bodyParser.json());
    app.use(cors(corsOptions));
    app.use("/api/v1/auth", authRouter);
    app.use("/api/v1/users", userRouter);
    app.use("/api/v1/users/:userId/tasks", taskRouter);
    app.use(globalErrorHandler);

    return app;
};
