import * as yup from "yup";


/**
 * Validation schema for user registration using Yup.
 * This schema validates the email, password, and confirmPassword fields for a registration form.
 * - `email`: Must be a valid email address and is required.
 * - `password`: Must be a string between 6 and 24 characters and is required.
 * - `confirmPassword`: Must match the password field and is required.
 * @type {yup.ObjectSchema}
 */
export const RegistrationSchema = yup.object(
    {
        email: yup.string()
                  .email("Invalid email address")
                  .required("Email is required"),
        password: yup.string()
                     .min(6, "Password must be at least 6 characters")
                     .max(24, "Password must be up to 24 characters")
                     .required("Password is required"),
        confirmPassword: yup.string()
                            .oneOf([yup.ref("password")], "Passwords must match")
                            .required("Confirm password is required")
    }
);
