import { Button, Dialog, DialogActions, DialogContent, DialogTitle, Typography } from "@mui/material";
import { FC }                                                                    from "react";


interface ConfirmationDialogProps {
    title: string;
    content: string;
    open: boolean;
    onConfirm: () => void;
    onCancel: () => void;
    onClose: () => void;
}

const ConfirmationDialog: FC<ConfirmationDialogProps> = ({ title, content, open, onConfirm, onCancel, onClose }) => {
    return (
        <Dialog open={ open } onClose={ onClose } maxWidth="sm" fullWidth>
            <DialogTitle>{ title }</DialogTitle>

            <DialogContent>
                <Typography variant="h6" gutterBottom>{ content }</Typography>
            </DialogContent>

            <DialogActions>
                <Button onClick={ onConfirm } color="error" variant="outlined">Confirm</Button>
                <Button onClick={ onCancel } color="secondary">Cancel</Button>
            </DialogActions>
        </Dialog>
    );
};

export default ConfirmationDialog;
