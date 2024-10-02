import { Alert, Box, Button, CircularProgress, Snackbar, TextField, Typography } from "@mui/material";
import { FC, useState }                                                          from "react";
import { useSelector }                                                           from "react-redux";
import {
    deleteUser,
    parseErrorMessage,
    RootState,
    updateUser,
    useAppDispatch,
    UserDeleteRequest,
    UserUpdateRequest
}                                                                                from "../../core";


const ProfilePage: FC = () => {
    const dispatch = useAppDispatch();
    const { userId, userEmail, appIsLoading, appError } = useSelector((state: RootState) => state.globals);

    const [currentPassword, setCurrentPassword] = useState("");
    const [newPassword, setNewPassword] = useState("");
    const [newPasswordConfirmation, setNewPasswordConfirmation] = useState("");
    const [deletePassword, setDeletePassword] = useState("");

    const [errorMessage, setErrorMessage] = useState<string | null>(null);
    const [showSnackbar, setShowSnackbar] = useState(false);

    const handleCloseSnackbar = () => {
        setShowSnackbar(false);
        setErrorMessage(null);
    };

    const validatePasswords = (): boolean => {
        if (!currentPassword || !newPassword || !newPasswordConfirmation) {
            setErrorMessage("All password fields are required.");
            setShowSnackbar(true);
            return false;
        }
        if (newPassword !== newPasswordConfirmation) {
            setErrorMessage("New password and confirmation password must match.");
            setShowSnackbar(true);
            return false;
        }
        return true;
    };

    const handleUpdatePassword = async () => {
        if (!validatePasswords()) {
            return;
        }

        const updReq: UserUpdateRequest = {
            userId: userId,
            userUpdateDTO: {
                currentPassword: currentPassword,
                newPassword: newPassword,
                newPasswordConfirmation: newPasswordConfirmation
            }
        };

        try {
            await dispatch(updateUser(updReq)).unwrap();
        } catch (error) {
            const errMsg = parseErrorMessage(error);
            setErrorMessage(errMsg);
            setShowSnackbar(true);
        }
    };

    const handleDeleteAccount = async () => {
        if (!deletePassword) {
            setErrorMessage("Password is required to delete account.");
            setShowSnackbar(true);
            return;
        }

        const delReq: UserDeleteRequest = {
            userId: userId,
            userDeletionDTO: {
                email: userEmail,
                currentPassword: deletePassword
            }
        };

        try {
            await dispatch(deleteUser(delReq)).unwrap();
        } catch (error) {
            const errMsg = parseErrorMessage(error);
            setErrorMessage(errMsg);
            setShowSnackbar(true);
        }
    };

    return (
        <Box sx={ { maxWidth: 400, mx: "auto", mt: 4 } }>
            <Typography variant="h5" gutterBottom>
                Account Settings
            </Typography>

            { appError && (
                <Alert severity="error" sx={ { mb: 2 } }>
                    { appError }
                </Alert>
            ) }

            <Snackbar
                open={ showSnackbar }
                autoHideDuration={ 6000 }
                onClose={ handleCloseSnackbar }
                message={ errorMessage }
                anchorOrigin={ { vertical: "top", horizontal: "center" } }
            />

            <Box component="section" sx={ { mb: 4 } }>
                <Typography variant="h6" gutterBottom>
                    Change Password
                </Typography>
                <TextField
                    label="Current Password"
                    type="password"
                    fullWidth
                    margin="normal"
                    value={ currentPassword }
                    onChange={ (e) => setCurrentPassword(e.target.value) }
                    required
                />
                <TextField
                    label="New Password"
                    type="password"
                    fullWidth
                    margin="normal"
                    value={ newPassword }
                    onChange={ (e) => setNewPassword(e.target.value) }
                    required
                />
                <TextField
                    label="Confirm New Password"
                    type="password"
                    fullWidth
                    margin="normal"
                    value={ newPasswordConfirmation }
                    onChange={ (e) => setNewPasswordConfirmation(e.target.value) }
                    required
                />
                <Button
                    variant="contained"
                    color="success"
                    fullWidth
                    onClick={ handleUpdatePassword }
                    disabled={ appIsLoading }
                >
                    { appIsLoading ? <CircularProgress size={ 24 }/> : "Update Password" }
                </Button>
            </Box>

            <Box component="section">
                <Typography variant="h6" gutterBottom>
                    Delete Account
                </Typography>
                <TextField
                    label="Email"
                    type="email"
                    fullWidth
                    margin="normal"
                    value={ userEmail }
                    disabled
                />
                <TextField
                    label="Current Password"
                    type="password"
                    fullWidth
                    margin="normal"
                    value={ deletePassword }
                    onChange={ (e) => setDeletePassword(e.target.value) }
                    required
                />
                <Button
                    variant="contained"
                    color="error"
                    fullWidth
                    onClick={ handleDeleteAccount }
                    disabled={ appIsLoading }
                >
                    { appIsLoading ? <CircularProgress size={ 24 }/> : "Delete Account" }
                </Button>
            </Box>
        </Box>
    );
};

export default ProfilePage;
