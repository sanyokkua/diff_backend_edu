import { CssBaseline }              from "@mui/material";
import { FC, JSX, useEffect }       from "react";
import { useNavigate }              from "react-router-dom";
import { LogLevel, useAppSelector } from "./core";
import { Layout }                   from "./react";


const logger = LogLevel.getLogger("App");

/**
 * The main App component that sets up the application layout and handles user authentication state.
 *
 * @component
 * @returns {JSX.Element} The rendered App component.
 */
const App: FC = (): JSX.Element => {
    const navigate = useNavigate();
    const isUserLoggedIn = useAppSelector((state) => state.users.userIsLoggedIn);

    /**
     * useEffect hook to check if the user is logged in and redirect to the login page if not.
     *
     * @param {boolean} isUserLoggedIn - Indicates if the user is logged in.
     * @param {Function} navigate - Function to navigate to different routes.
     */
    useEffect(() => {
        if (!isUserLoggedIn) {
            logger.info("User is not logged in, redirecting to login page.");
            navigate("/login");
        }
    }, [isUserLoggedIn, navigate]);

    logger.debug("Rendering App component");

    return (
        <>
            <CssBaseline/>
            <Layout/>
        </>
    );
};

export default App;
