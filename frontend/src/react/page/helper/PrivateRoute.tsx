import { FC, JSX }                  from "react";
import { Navigate, Outlet }         from "react-router-dom";
import { LogLevel, useAppSelector } from "../../../core";


const logger = LogLevel.getLogger("PrivateRoute");

/**
 * A component that restricts access to its children based on user authentication status.
 * @function PrivateRoute
 * @returns {JSX.Element} - The rendered private route component.
 */
const PrivateRoute: FC = (): JSX.Element => {
    logger.debug("Private route render");
    const { userIsLoggedIn } = useAppSelector((state) => state.users);
    return userIsLoggedIn ? <Outlet/> : <Navigate to="/"/>;
};

export default PrivateRoute;
