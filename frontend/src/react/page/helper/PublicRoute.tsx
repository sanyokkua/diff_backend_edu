import { FC, JSX }                  from "react";
import { Navigate, Outlet }         from "react-router-dom";
import { LogLevel, useAppSelector } from "../../../core";


const logger = LogLevel.getLogger("PublicRoute");

/**
 * A component that restricts access to its children based on user authentication status.
 * @function PublicRoute
 * @returns {JSX.Element} - The rendered public route component.
 */
const PublicRoute: FC = (): JSX.Element => {
    logger.debug("Public route render");
    const { userIsLoggedIn } = useAppSelector((state) => state.users);
    return userIsLoggedIn ? <Navigate to="/dashboard"/> : <Outlet/>;
};

export default PublicRoute;
