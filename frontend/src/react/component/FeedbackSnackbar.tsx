import { Snackbar } from "@mui/material";
import Alert        from "@mui/material/Alert";
import { FC }       from "react";
import { LogLevel } from "../../core";

// Initialize logger for the FeedbackSnackbar component
const logger = LogLevel.getLogger("FeedbackSnackbar");

/**
 * FeedbackType - Defines the structure for feedback messages.
 * @property message - The content of the feedback message.
 * @property severity - The severity level of the feedback, either "error" or "success".
 */
export type FeedbackType = {
    message: string;
    severity: "error" | "success";
};

/**
 * Props for the FeedbackSnackbar component.
 * @property feedbackMessage - The message to be displayed in the Snackbar, containing the message and severity.
 * @property onClose - Callback function triggered when the Snackbar is closed.
 */
interface FeedbackSnackbarProps {
    feedbackMessage: FeedbackType | null | undefined;
    onClose: () => void;
}

/**
 * FeedbackSnackbar Component
 * Renders a Snackbar that shows success or error messages to the user.
 * Automatically hides after 6 seconds or when the user closes it.
 *
 * @param feedbackMessage - The message and its severity to be displayed in the Snackbar.
 * @param onClose - Function to handle the closing of the Snackbar.
 */
const FeedbackSnackbar: FC<FeedbackSnackbarProps> = ({ feedbackMessage, onClose }) => {
    // Check if there's a feedback message to display
    if (!feedbackMessage) {
        logger.debug("No feedback message to display."); // Debug log for when no feedback message is available
        return <></>; // Return nothing if no message is present
    }

    logger.info(`Displaying ${ feedbackMessage.severity } feedback: "${ feedbackMessage.message }"`); // Log the feedback message and its severity

    // Render the Snackbar with the feedback message and severity
    return (
        <Snackbar open={ !!feedbackMessage } autoHideDuration={ 6000 } onClose={ onClose }>
            <Alert onClose={ onClose } severity={ feedbackMessage.severity }>
                { feedbackMessage.message }
            </Alert>
        </Snackbar>
    );
};

export default FeedbackSnackbar;
