import { Box, Button, CircularProgress, TextField, Typography } from "@mui/material";
import { FC, JSX, useEffect, useState }                         from "react";
import { useSelector }                                          from "react-redux";
import {
    deleteUser,
    LogLevel,
    parseErrorMessage,
    RootState,
    setFeedback,
    setHeaderTitle,
    updateUser,
    useAppDispatch,
    UserDeleteRequest,
    UserUpdateRequest
}                                                               from "../../core";


const logger = LogLevel.getLogger("ProfilePage");
/**
 * ProfilePage component is a functional React component that provides users with the ability to
 * update their password and delete their account. It interacts with the Redux store to manage
 * user state and provides feedback through the UI.
 *
 * @component
 * @returns {JSX.Element} The rendered ProfilePage component.
 */
const ProfilePage: FC = (): JSX.Element => {
    const dispatch = useAppDispatch();
    const { userId, userEmail, userIsLoading } = useSelector((state: RootState) => state.users);

    // State hooks for form inputs and feedback
    const [currentPassword, setCurrentPassword] = useState<string>("");
    const [newPassword, setNewPassword] = useState<string>("");
    const [newPasswordConfirmation, setNewPasswordConfirmation] = useState<string>("");
    const [deletePassword, setDeletePassword] = useState<string>("");

    // Set the header title on component mount
    useEffect(() => {
        dispatch(setHeaderTitle("Profile"));
    }, [dispatch]);

    /**
     * Validates the password fields for update functionality.
     * Checks if all fields are filled and if the new password matches the confirmation.
     *
     * @returns {boolean} True if passwords are valid; otherwise, false.
     */
    const validatePasswords = (): boolean => {
        if (!currentPassword || !newPassword || !newPasswordConfirmation) {
            dispatch(setFeedback({ message: "All password fields are required.", severity: "error" }));
            return false;
        }
        if (newPassword !== newPasswordConfirmation) {
            dispatch(setFeedback({ message: "New password and confirmation must match.", severity: "error" }));
            return false;
        }
        return true;
    };

    /**
     * Handles the password update process.
     * Validates input, constructs a UserUpdateRequest object, and dispatches the update action.
     * Provides feedback on success or failure.
     *
     * @async
     * @returns {Promise<void>} A promise that resolves when the update is complete.
     */
    const handleUpdatePassword = async (): Promise<void> => {
        logger.debug("Attempting to update password");
        if (!validatePasswords()) {
            return;
        }

        const updateRequest: UserUpdateRequest = {
            userId,
            userUpdateDTO: {
                currentPassword,
                newPassword,
                newPasswordConfirmation
            }
        };

        try {
            await dispatch(updateUser(updateRequest)).unwrap();
            dispatch(setFeedback({ message: "Password updated successfully.", severity: "success" }));
        } catch (error) {
            const errorMessage = parseErrorMessage(error, "Failed to update password");
            dispatch(setFeedback({ message: errorMessage, severity: "error" }));
            logger.error(errorMessage);
        }
    };

    /**
     * Handles the account deletion process.
     * Validates the input, constructs a UserDeleteRequest object, and dispatches the delete action.
     * Provides feedback on success or failure.
     *
     * @async
     * @returns {Promise<void>} A promise that resolves when the deletion is complete.
     */
    const handleDeleteAccount = async (): Promise<void> => {
        logger.debug("Attempting to delete account");
        if (!deletePassword) {
            dispatch(setFeedback({ message: "Password is required to delete the account.", severity: "error" }));
            return;
        }

        const deleteRequest: UserDeleteRequest = {
            userId,
            userDeletionDTO: {
                email: userEmail,
                currentPassword: deletePassword
            }
        };

        try {
            await dispatch(deleteUser(deleteRequest)).unwrap();
            dispatch(setFeedback({ message: "Account deleted successfully.", severity: "success" }));
        } catch (error) {
            const errorMessage = parseErrorMessage(error, "Failed to delete account");
            dispatch(setFeedback({ message: errorMessage, severity: "error" }));
            logger.error(errorMessage);
        }
    };

    return (
        <>
            { userIsLoading && <CircularProgress size={ 24 }/> }
            <Box sx={ { maxWidth: 400, mx: "auto", mt: 4 } }>
                <Typography variant="h5" gutterBottom>
                    Account Settings
                </Typography>

                <Box component="section" sx={ { mb: 4 } }>
                    <Typography variant="h6" gutterBottom>
                        Change Password
                    </Typography>
                    <TextField
                        label="Current Password"
                        type="password"
                        margin="normal"
                        value={ currentPassword }
                        fullWidth
                        required
                        onChange={ (e) => setCurrentPassword(e.target.value) }
                    />
                    <TextField
                        label="New Password"
                        type="password"
                        margin="normal"
                        value={ newPassword }
                        fullWidth
                        required
                        onChange={ (e) => setNewPassword(e.target.value) }
                    />
                    <TextField
                        label="Confirm New Password"
                        type="password"
                        margin="normal"
                        value={ newPasswordConfirmation }
                        fullWidth
                        required
                        onChange={ (e) => setNewPasswordConfirmation(e.target.value) }
                    />
                    <Button
                        variant="contained"
                        color="success"
                        fullWidth
                        onClick={ handleUpdatePassword }
                        disabled={ userIsLoading }
                    >
                        Update Password
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
                        margin="normal"
                        fullWidth
                        required
                        value={ deletePassword }
                        onChange={ (e) => setDeletePassword(e.target.value) }
                    />
                    <Button
                        variant="contained"
                        color="error"
                        fullWidth
                        disabled={ userIsLoading }
                        onClick={ handleDeleteAccount }
                    >
                        Delete Account
                    </Button>
                </Box>
            </Box>
        </>
    );
};

export default ProfilePage;
