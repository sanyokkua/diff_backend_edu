import { Box }    from "@mui/material";
import React      from "react";
import { Outlet } from "react-router-dom";
import AppToolbar from "./AppToolbar";


const Layout: React.FC = (): React.JSX.Element => {
    return (
        <Box>
            <AppToolbar/>
            <Outlet/>
        </Box>
    );
};

export default Layout;