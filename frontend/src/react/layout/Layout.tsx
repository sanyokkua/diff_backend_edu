import { Box, CircularProgress, SxProps }                          from "@mui/material";
import { FC, JSX }                                                 from "react";
import { Outlet }                                                  from "react-router-dom";
import { clearFeedback, LogLevel, useAppDispatch, useAppSelector } from "../../core";
import FeedbackSnackbar                                            from "../component/FeedbackSnackbar";
import AppToolbar                                                  from "./AppToolbar";


const logger = LogLevel.getLogger("Layout");

const layoutStyles: SxProps = {
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "center",
    width: "80%",
    margin: "0 auto",
    padding: 2
};

/**
 * Layout component for the application.
 * @function Layout
 * @returns {JSX.Element} - The rendered layout component.
 */
const Layout: FC = (): JSX.Element => {
    const dispatch = useAppDispatch();
    const { feedback, isLoading } = useAppSelector((state) => state.globals);

    /**
     * Handles the closing of the feedback notification.
     * @function handleNotificationClose
     * @returns {void}
     */
    const handleNotificationClose = (): void => {
        logger.debug("Clearing Global Notification");
        dispatch(clearFeedback());
    };

    return (
        <>
            { isLoading ? (
                <CircularProgress/>
            ) : (
                  <Box>
                      <AppToolbar/>
                      <Box sx={ layoutStyles }>
                          <Outlet/>
                      </Box>
                  </Box>
              ) }
            { feedback.message && <FeedbackSnackbar feedbackMessage={ feedback } onClose={ handleNotificationClose }/> }
        </>
    );
};

export default Layout;
