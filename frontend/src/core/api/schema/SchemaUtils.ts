import * as yup         from "yup";
import { FormResponse } from "../../../react";


/**
 * Validates form data against a given Yup schema and handles validation results.
 * @param {FormResponse} formData - The form data to be validated.
 * @param {yup.AnyObject} currentSchema - The Yup schema to validate against.
 * @param {function(Partial<FormResponse>): void} onValidationFinished - Callback function to handle validation results.
 * @returns {Promise<boolean>} - A promise that resolves to true if validation passes, otherwise false.
 */
export const validateForm = async (
    formData: FormResponse,
    currentSchema: yup.AnyObject,
    onValidationFinished: (result: Partial<FormResponse>) => void
): Promise<boolean> => {
    try {
        await currentSchema.validate(formData, { abortEarly: false });
        onValidationFinished({});
        return true;
    } catch (validationErrors) {
        const formErrors: Partial<FormResponse> = {};
        if (validationErrors instanceof yup.ValidationError) {
            validationErrors.inner.forEach((err) => {
                if (err.path) {
                    formErrors[err.path as keyof FormResponse] = err.message;
                }
            });
        }
        onValidationFinished(formErrors);
        return false;
    }
};
