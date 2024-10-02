import { Box, Button, Typography }                                 from "@mui/material";
import { ChangeEvent, FC, FormEvent, useState }                    from "react";
import { Link }                                                    from "react-router-dom";
import { LoginSchema, LogLevel, RegistrationSchema, validateForm } from "../../core";
import FeedbackSnackbar, { FeedbackType }                          from "./FeedbackSnackbar";
import FormField                                                   from "./FormField";


const log = LogLevel.getLogger("FormGeneric");

// Types for form data response and form types
export type FormResponse = {
    email: string;
    password: string;
    confirmPassword?: string;
}; // Form fields structure, confirmPassword is only for Registration

type FormType = "Login" | "Register"; // Two types of forms

type FormConfiguration = {
    title: string;
    mainBtn: string;
    altText: string;
    altLink: string;
}; // Configuration for the UI based on form type

// Form configuration based on the type of form (Login/Register)
const formConfigs: Record<FormType, FormConfiguration> = {
    Login: {
        title: "Login to Existing Account",
        mainBtn: "Login",
        altText: "Don't have an account? Register",
        altLink: "/register"
    },
    Register: {
        title: "Create New Account",
        mainBtn: "Register",
        altText: "Already have an account? Login",
        altLink: "/login"
    }
};

// Initial values for the form inputs
const initialFormValues: FormResponse = {
    email: "",
    password: "",
    confirmPassword: "" // Only needed for Register form
};

// Props for the form component
interface FormGenericProps {
    formType: FormType; // The form type ("Login" or "Register")
    onSubmit: (data: FormResponse) => Promise<void>; // Callback function triggered on form submission
}

/**
 * FormGeneric Component
 * This component renders either a Login or Registration form based on the formType prop.
 * It handles form state, validation, and submission.
 */
const FormGeneric: FC<FormGenericProps> = ({ formType, onSubmit }) => {
    log.info(`Rendering FormGeneric component for ${ formType } form.`); // Log rendering of form component

    // Select appropriate schema (Login or Registration) for validation
    const currentSchema = formType === "Login" ? LoginSchema : RegistrationSchema;

    // Destructure UI configuration for the selected form type
    const { title, mainBtn, altText, altLink } = formConfigs[formType];

    // State variables to manage form data, validation errors, and feedback messages
    const [formData, setFormData] = useState<FormResponse>(initialFormValues); // Input values
    const [errors, setErrors] = useState<Partial<FormResponse>>({});           // Validation errors
    const [feedbackMessage, setFeedbackMessage] = useState<FeedbackType | null>(null); // Feedback snackbar message

    /**
     * Handle input field changes.
     * This updates the formData state dynamically when any form input is changed.
     */
    const handleInputChange = (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { name, value } = e.target;
        log.debug(`Input changed: ${ name }, value length: ${ value.length }`); // Log input changes without logging actual value
        setFormData((prev) => ({ ...prev, [name]: value }));
    };

    /**
     * Handle form submission.
     * This function validates the form, calls the onSubmit callback, and shows success or error feedback.
     */
    const handleFormSubmit = async (e: FormEvent) => {
        e.preventDefault(); // Prevent default form submission behavior (page reload)
        log.debug("Form submission started.");

        // Validate form data using the appropriate schema
        const isValid = await validateForm(formData, currentSchema, setErrors);
        if (isValid) {
            log.info("Form validation passed.");
            // If the form is valid, try submitting the form data
            try {
                log.debug("Attempting to submit form data.");
                await onSubmit(formData); // Trigger form submission logic passed as a prop
                setFeedbackMessage({ message: "Form submitted successfully", severity: "success" }); // Show success feedback
                log.info("Form submitted successfully.");
            } catch (error) {
                log.error("Form submission failed.", error); // Log error details
                let errMsg: string = "Submission failed. Please try again later.";
                if (typeof error === "string") {
                    errMsg = error;
                } else if (error instanceof Error) {
                    errMsg = error.message;
                }
                setFeedbackMessage({ message: errMsg, severity: "error" });
            }
        } else {
            log.warn("Form validation failed.", { errors }); // Log validation errors
            setFeedbackMessage({ message: "Form validation failed", severity: "error" });
        }
    };

    /**
     * Handle closing of the feedback snackbar.
     * This clears the feedbackMessage state, closing the snackbar.
     */
    const handleCloseSnackbar = () => setFeedbackMessage(null);

    return (
        <Box display="flex" flexDirection="column" alignItems="center" justifyContent="center" minHeight="100vh"
             p={ 2 }>
            <Typography variant="h4" gutterBottom>
                { title } {/* Form title (Login/Register) */ }
            </Typography>

            <Box display="flex" flexDirection="column" alignItems="center" width="100%" maxWidth="400px">
                {/* Form element with submit handler */ }
                <form onSubmit={ handleFormSubmit } style={ { width: "100%" } }>
                    {/* Email field */ }
                    <FormField
                        fieldName={ "email" }
                        label="Email"
                        type="email"
                        value={ formData.email }
                        error={ errors.email }
                        onChange={ handleInputChange }
                    />
                    {/* Password field */ }
                    <FormField
                        fieldName={ "password" }
                        label="Password"
                        type="password"
                        value={ formData.password }
                        error={ errors.password }
                        onChange={ handleInputChange }
                    />
                    {/* Confirm password field (only shown for Register form) */ }
                    { formType === "Register" && (
                        <FormField
                            fieldName={ "confirmPassword" }
                            label="Confirm Password"
                            type="password"
                            value={ formData.confirmPassword ?? "" }
                            error={ errors.confirmPassword }
                            onChange={ handleInputChange }
                        />
                    ) }
                    {/* Main button (Login/Register) */ }
                    <Button variant="contained" color="primary" fullWidth type="submit">
                        { mainBtn }
                    </Button>
                </form>

                {/* Link to alternate form (e.g., switch between Login/Register) */ }
                <Button component={ Link } to={ altLink } variant="text" color="secondary" fullWidth sx={ { mt: 2 } }>
                    { altText }
                </Button>
            </Box>

            {/* Feedback snackbar for showing success/error messages */ }
            <FeedbackSnackbar feedbackMessage={ feedbackMessage } onClose={ handleCloseSnackbar }/>
        </Box>
    );
};

export default FormGeneric;
