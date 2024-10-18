import { DataSource } from "typeorm";
import { Task, User } from "../model";


export const getDataSource = (): DataSource => {
    return new DataSource(
        {
            type: "postgres",
            host: process.env.DB_HOST,
            port: Number(process.env.DB_PORT),
            username: process.env.DB_USERNAME,
            password: process.env.DB_PASSWORD,
            database: process.env.DB_NAME,
            entities: [User, Task],
            synchronize: true
        }
    );
};