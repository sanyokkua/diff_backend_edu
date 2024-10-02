import { AppBar, Box, Toolbar, Typography } from "@mui/material";
import { FC }                               from "react";
import { Link }                             from "react-router-dom";
import { useAppSelector }                   from "../../core";


const AppToolbar: FC = () => {
    const { appBarHeader, userIsLoggedIn, userEmail } = useAppSelector((state) => state.globals);

    return (
        <AppBar position="sticky">
            <Toolbar>
                <Box sx={ { display: "flex", justifyContent: "space-between", width: "100%" } }>
                    <Typography variant="h6" component="div">
                        { userIsLoggedIn ? appBarHeader : "Task Manager" }
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
