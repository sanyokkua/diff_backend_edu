import { CssBaseline }      from "@mui/material";
import React, { useEffect } from "react";
import { useNavigate }      from "react-router-dom";
import { useAppSelector }   from "./core";
import { Layout }           from "./react";


const App: React.FC = (): React.JSX.Element => {
    const navigate = useNavigate();
    const { userIsLoggedIn } = useAppSelector((state) => state.globals);

    useEffect(() => {
        if (!userIsLoggedIn) {
            navigate("/login");
        }
    }, [navigate, userIsLoggedIn]);

    return (
        <>
            <CssBaseline/>
            <Layout/>
        </>
    );
};

export default App;