import { Navigate, Outlet } from "react-router-dom";
import { useAppSelector }   from "../../../core";


const PublicRoute = () => {
    const { userIsLoggedIn } = useAppSelector((state) => state.globals);
    return userIsLoggedIn ? <Navigate to="/dashboard"/> : <Outlet/>;
};

export default PublicRoute;
