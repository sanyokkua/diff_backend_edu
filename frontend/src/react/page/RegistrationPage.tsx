import { FC, JSX, useEffect }                                                                   from "react";
import { useNavigate }                                                                          from "react-router-dom";
import { LogLevel, registerUser, setFeedback, setHeaderTitle, useAppDispatch, UserCreationDTO } from "../../core";
import { FormGeneric, FormResponse }                                                            from "../component";


const logger = LogLevel.getLogger("RegistrationPage");

/**
 * RegistrationPage component.
 *
 * This component is responsible for rendering the registration page and handling user registration.
 * It sets the header title to "Registration" and manages the form submission process.
 *
 * @component
 * @returns {JSX.Element} The rendered registration page component.
 */
const RegistrationPage: FC = (): JSX.Element => {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

    useEffect(() => {
        // Set the header title to "Registration" when the component mounts.
        dispatch(setHeaderTitle("Registration"));
    }, [dispatch]);

    /**
     * Handles the form submission.
     *
     * @param {FormResponse} formData - The form data containing email, password, and password confirmation.
     */
    const handleFormSubmit = async (formData: FormResponse) => {
        logger.debug("Form submitted", formData);
        const userCreationData: UserCreationDTO = {
            email: formData.email,
            password: formData.password,
            passwordConfirmation: formData.confirmPassword ?? ""
        };

        try {
            await dispatch(registerUser(userCreationData)).unwrap();
            dispatch(setFeedback({ message: "Registration successful", severity: "success" }));
            setTimeout(() => navigate("/dashboard"), 2000);
        } catch (error) {
            logger.error("Registration error:", error);
            dispatch(setFeedback({ message: "Registration failed. Please try again.", severity: "error" }));
        }
    };

    return <FormGeneric formType="Register" onSubmit={ handleFormSubmit }/>;
};

export default RegistrationPage;
