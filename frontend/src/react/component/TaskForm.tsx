import { Box, Button, TextareaAutosize, TextField, Typography } from "@mui/material";
import React, { useState }                                      from "react";
import { TaskDto }                                              from "../../core";


interface TaskFormProps {
    task: TaskDto;
    onSave: (task: TaskDto) => Promise<void>;
    onCancel: () => void;
    error: string;
}

const TaskForm: React.FC<TaskFormProps> = ({ task, onSave, onCancel, error }) => {
    const [name, setName] = useState(task.name);
    const [description, setDescription] = useState(task.description);
    const [nameError, setNameError] = useState("");

    const handleSave = () => {
        if (name.trim() === "") {
            setNameError("Name cannot be empty");
            return;
        }
        setNameError("");
        onSave({ ...task, name, description });
    };

    return (
        <Box
            component="form"
            sx={ {
                display: "flex",
                flexDirection: "column",
                gap: 2,
                width: "80%",
                margin: "0 auto",
                padding: 2
            } }
        >
            <TextField
                label="Name"
                value={ name }
                onChange={ (e) => setName(e.target.value) }
                variant="outlined"
                error={ !!nameError }
                helperText={ nameError }
            />
            <TextareaAutosize
                minRows={ 4 }
                placeholder="Description"
                value={ description }
                onChange={ (e) => setDescription(e.target.value) }
                style={ { width: "100%", padding: "8px", fontSize: "16px" } }
            />

            { error && (
                <Typography color="error" variant="body2">
                    { error }
                </Typography>
            ) }

            <Box sx={ { display: "flex", justifyContent: "flex-end", gap: 1 } }>
                <Button variant="outlined" color="secondary" onClick={ onCancel }>
                    Cancel
                </Button>
                <Button variant="contained" color="primary" onClick={ handleSave }>
                    Save
                </Button>
            </Box>
        </Box>
    );
};

export default TaskForm;
