import { createBrowserRouter }                                                       from "react-router-dom";
import App                                                                           from "../../App";
import { Dashboard, ErrorPage, Login, PrivateRoute, Profile, PublicRoute, Register } from "../page";
import TaskCreatePage                                                                from "../page/TaskCreatePage";
import TaskEditPage                                                                  from "../page/TaskEditPage";


const BrowserRouter = createBrowserRouter(
    [
        {
            path: "/",
            element: <App/>,
            errorElement: <ErrorPage/>,
            children: [
                {
                    element: <PublicRoute/>,
                    children: [
                        {
                            index: true,
                            element: <Login/>,
                            errorElement: <ErrorPage/>
                        }
                    ]
                },
                {
                    path: "/login",
                    element: <PublicRoute/>,
                    children: [
                        {
                            index: true,
                            element: <Login/>,
                            errorElement: <ErrorPage/>
                        }
                    ]
                },
                {
                    path: "/register",
                    element: <PublicRoute/>,
                    children: [
                        {
                            index: true,
                            element: <Register/>,
                            errorElement: <ErrorPage/>
                        }
                    ]
                },
                {
                    path: "/dashboard",
                    element: <PrivateRoute/>,
                    children: [
                        {
                            index: true,
                            element: <Dashboard/>,
                            errorElement: <ErrorPage/>
                        },
                        {
                            path: "/dashboard/new",
                            element: <TaskCreatePage/>,
                            errorElement: <ErrorPage/>
                        },
                        {
                            path: "/dashboard/edit",
                            element: <TaskEditPage/>,
                            errorElement: <ErrorPage/>
                        }
                    ]
                },
                {
                    path: "/profile",
                    element: <PrivateRoute/>,
                    children: [
                        {
                            index: true,
                            element: <Profile/>,
                            errorElement: <ErrorPage/>
                        }
                    ]
                }
            ]
        }
    ]);

export default BrowserRouter;
