import { FC, JSX, useEffect }                                                             from "react";
import { useNavigate }                                                                    from "react-router-dom";
import { loginUser, LogLevel, setFeedback, setHeaderTitle, useAppDispatch, UserLoginDto } from "../../core";
import { FormGeneric, FormResponse }                                                      from "../component";


const logger = LogLevel.getLogger("LoginPage");

/**
 * LoginPage component.
 *
 * This component is responsible for rendering the login page and handling user login.
 * It sets the header title to "Login" and manages the form submission process.
 *
 * @component
 * @returns {JSX.Element} The rendered login page component.
 */
const LoginPage: FC = (): JSX.Element => {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

    useEffect(() => {
        // Set the header title to "Login" when the component mounts.
        dispatch(setHeaderTitle("Login"));
    }, [dispatch]);

    /**
     * Handles the form submission.
     *
     * @param {FormResponse} formData - The form data containing email and password.
     */
    const handleFormSubmit = async (formData: FormResponse) => {
        logger.debug("Form submission initiated", formData);
        const loginRequest: UserLoginDto = {
            email: formData.email,
            password: formData.password
        };

        try {
            await dispatch(loginUser(loginRequest)).unwrap();
            dispatch(setFeedback({ message: "Login successful", severity: "success" }));
            setTimeout(() => navigate("/dashboard"), 2000);
        } catch (error) {
            logger.error("Login failed", error);
            dispatch(setFeedback({ message: "Login failed. Please try again.", severity: "error" }));
        }
    };

    return <FormGeneric formType="Login" onSubmit={ handleFormSubmit }/>;
};

export default LoginPage;
