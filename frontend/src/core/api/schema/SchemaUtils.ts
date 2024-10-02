import * as yup         from "yup";
import { FormResponse } from "../../../react/component/FormGeneric";


export const validateForm = async (formData: FormResponse, currentSchema: yup.AnyObject, onValidationFinished: (result: Partial<FormResponse>) => void): Promise<boolean> => {
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