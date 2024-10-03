import { AppBar, Box, Toolbar, Typography } from "@mui/material";
import { FC, JSX }                          from "react";
import { Link }                             from "react-router-dom";
import { LogLevel, useAppSelector }         from "../../core";


const logger = LogLevel.getLogger("AppToolbar");

/**
 * A toolbar component for the application.
 * @function AppToolbar
 * @returns {JSX.Element} - The rendered app toolbar.
 */
const AppToolbar: FC = (): JSX.Element => {
    const { headerTitle } = useAppSelector((state) => state.globals);
    const { userIsLoggedIn, userEmail } = useAppSelector((state) => state.users);

    logger.debug("Rendering AppToolbar");

    if (!headerTitle || !userEmail) {
        logger.warn("Missing required state properties: headerTitle or email");
    }

    return (
        <AppBar position="sticky">
            <Toolbar>
                <Box sx={ { display: "flex", justifyContent: "space-between", width: "100%" } }>
                    <Typography variant="h6" component="div">
                        { userIsLoggedIn ? headerTitle : "Task Manager" }
                    </Typography>

                    { userIsLoggedIn && (
                        <Link to="/profile" style={ { textDecoration: "none", color: "inherit" } }>
                            <Typography variant="h6" component="div">
                                { userEmail }
                            </Typography>
                        </Link>
                    ) }
                </Box>
            </Toolbar>
        </AppBar>
    );
};

export default AppToolbar;
