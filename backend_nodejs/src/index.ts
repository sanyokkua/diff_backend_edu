import "reflect-metadata";
import dotenv               from "dotenv";
import { createExpressApp } from "./app";


dotenv.config();
const port = process.env.PORT ?? 3000;

const App = createExpressApp();

App.listen(port, () => {
    console.log(`Server is running on port ${ port }`);
});