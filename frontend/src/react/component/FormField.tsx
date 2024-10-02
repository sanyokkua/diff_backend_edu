import { FormControl, TextField } from "@mui/material";
import { ChangeEvent, FC }        from "react";
import { LogLevel }               from "../../core";

// Initialize logger for the FormField component
const log = LogLevel.getLogger("FormField");

/**
 * FormFieldProps - Defines the properties for the FormField component.
 * @property fieldName - The name attribute for the input field, used to identify the field in forms.
 * @property label - The label displayed alongside the input field.
 * @property type - The input type (e.g., "text", "email", "password").
 * @property value - The current value of the input field.
 * @property error - Error message to display if validation fails, undefined if there's no error.
 * @property onChange - Callback function triggered when the input value changes.
 */
interface FormFieldProps {
    fieldName: string;
    label: string;
    type: string;
    value: string;
    error: string | undefined;
    onChange: (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void;
}

/**
 * FormField Component
 * A reusable form input field with built-in error handling and MUI styling.
 * Logs the field's state changes for debugging and development purposes.
 *
 * @param fieldName - The name of the field.
 * @param label - The label displayed for the field.
 * @param type - The type of input (e.g., "text", "password").
 * @param value - The current value of the field.
 * @param error - Error message (if any) to display for validation failures.
 * @param onChange - Event handler for field value changes.
 */
const FormField: FC<FormFieldProps> = ({ fieldName, label, type, value, error, onChange }) => {
    log.debug(`Rendering FormField: ${ fieldName }, type: ${ type }`); // Debug log for rendering the field

    return (
        <FormControl fullWidth margin="normal">
            <TextField
                name={ fieldName }
                label={ label }
                type={ type }
                value={ value }
                onChange={ (e) => {
                    log.debug(`Field "${ fieldName }" changed. New value length: ${ e.target.value.length }`); // Log value length change (not actual value for privacy)
                    onChange(e); // Trigger parent onChange handler
                } }
                error={ !!error }
                helperText={ error }
            />
        </FormControl>
    );
};

export default FormField;
