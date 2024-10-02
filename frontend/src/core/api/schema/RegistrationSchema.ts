import * as yup from "yup";


export const RegistrationSchema = yup.object(
    {
        email: yup.string().email("Invalid email address").required("Email is required"),
        password: yup
            .string()
            .min(6, "Password must be at least 6 characters")
            .max(24, "Password must be up to 24 characters")
            .required("Password is required"),
        confirmPassword: yup
            .string()
            .oneOf([yup.ref("password")], "Passwords must match")
            .required("Confirm password is required")
    }
);