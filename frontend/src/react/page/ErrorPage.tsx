import { Box, Container, Paper, Typography }        from "@mui/material";
import { FC, JSX, useEffect }                       from "react";
import { useRouteError }                            from "react-router-dom";
import { LogLevel, setHeaderTitle, useAppDispatch } from "../../core";


const logger = LogLevel.getLogger("ErrorPage");

/**
 * ErrorPage component.
 *
 * This component is responsible for displaying an error page when an unexpected error occurs.
 * It uses Redux to set the header title and logs the rendering process.
 *
 * @component
 * @returns {JSX.Element} The rendered error page component.
 */
const ErrorPage: FC = (): JSX.Element => {
    const dispatch = useAppDispatch();

    useEffect(() => {
        // Dispatch an action to set the header title to "Error" when the component mounts.
        dispatch(setHeaderTitle("Error"));
    }, [dispatch]);

    // Retrieve the error object from the route.
    const error = useRouteError();
    const errorDetails = error as { statusText?: string; message?: string };

    // Log the rendering of the ErrorPage component.
    logger.debug("Rendering ErrorPage component");

    return (
        <Container maxWidth="sm">
            <Paper elevation={ 3 } sx={ { padding: 4, marginTop: 4 } }>
                <Box textAlign="center">
                    <Typography variant="h1" component="h1" gutterBottom>
                        Oops!
                    </Typography>
                    <Typography variant="body1" gutterBottom>
                        Sorry, an unexpected error has occurred.
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                        <i>{ errorDetails?.statusText ?? errorDetails?.message ?? "Unknown error" }</i>
                    </Typography>
                </Box>
            </Paper>
        </Container>
    );
};

export default ErrorPage;
