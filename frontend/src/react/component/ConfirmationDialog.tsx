import { Button, Dialog, DialogActions, DialogContent, DialogTitle, Typography } from "@mui/material";
import { FC, JSX }                                                               from "react";
import { LogLevel }                                                              from "../../core";


const logger = LogLevel.getLogger("ConfirmationDialog");

/**
 * Props for the ConfirmationDialog component.
 * @property {string} title - The title of the dialog.
 * @property {string} message - The message displayed in the dialog.
 * @property {boolean} isOpen - Whether the dialog is open.
 * @property {Function} onConfirm - Function to call when the confirm button is clicked.
 * @property {Function} onCancel - Function to call when the cancel button is clicked.
 * @property {Function} onClose - Function to call when the dialog is closed.
 */
interface ConfirmationDialogProps {
    title: string;
    message: string;
    isOpen: boolean;
    onConfirm: () => void;
    onCancel: () => void;
    onClose: () => void;
}

/**
 * A confirmation dialog component.
 * @function ConfirmationDialog
 * @param props - The props for the component.
 * @returns {JSX.Element} - The rendered confirmation dialog.
 */
const ConfirmationDialog: FC<ConfirmationDialogProps> = ({
                                                             title,
                                                             message,
                                                             isOpen,
                                                             onConfirm,
                                                             onCancel,
                                                             onClose
                                                         }): JSX.Element => {
    logger.debug("Rendering ConfirmationDialog");

    const handleConfirm = () => {
        try {
            onConfirm();
            logger.info("User confirmed the action");
        } catch (error) {
            logger.error("Error during confirmation", error);
        }
    };

    const handleCancel = () => {
        try {
            onCancel();
            logger.info("User canceled the action");
        } catch (error) {
            logger.error("Error during cancellation", error);
        }
    };

    return (
        <Dialog open={ isOpen } onClose={ onClose } maxWidth="sm" fullWidth>
            <DialogTitle>{ title }</DialogTitle>
            <DialogContent>
                <Typography variant="h6" gutterBottom>{ message }</Typography>
            </DialogContent>
            <DialogActions>
                <Button onClick={ handleConfirm } color="error" variant="outlined">Confirm</Button>
                <Button onClick={ handleCancel } color="secondary">Cancel</Button>
            </DialogActions>
        </Dialog>
    );
};

export default ConfirmationDialog;
