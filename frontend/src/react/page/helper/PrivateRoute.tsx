import { Navigate, Outlet } from "react-router-dom";
import { useAppSelector }   from "../../../core";


const PrivateRoute = () => {
    const { userIsLoggedIn } = useAppSelector((state) => state.globals);
    return userIsLoggedIn ? <Outlet/> : <Navigate to="/"/>;
};

export default PrivateRoute;