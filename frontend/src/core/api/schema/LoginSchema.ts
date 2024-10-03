import * as yup from "yup";


/**
 * Validation schema for user login using Yup.
 * This schema validates the email and password fields for a login form.
 * - `email`: Must be a valid email address and is required.
 * - `password`: Must be a non-empty string and is required.
 * @type {yup.ObjectSchema}
 */
export const LoginSchema = yup.object(
    {
        email: yup.string()
                  .email("Invalid email address")
                  .required("Email is required"),
        password: yup.string()
                     .required("Password is required")
    }
);
