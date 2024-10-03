import { Box, Button, TextareaAutosize, TextField, Typography } from "@mui/material";
import React, { JSX, useState }                                 from "react";
import { LogLevel, setFeedback, TaskDto, useAppDispatch }       from "../../core";


const logger = LogLevel.getLogger("TaskForm");

/**
 * Props for the TaskForm component.
 * @property {TaskDto} task - The task data.
 * @property {Function} onSave - Function to call when saving the task.
 * @property {Function} onCancel - Function to call when cancelling the form.
 * @property {string} errorMessage - Error message to display.
 */
interface TaskFormProps {
    task: TaskDto;
    onSave: (task: TaskDto) => Promise<void>;
    onCancel: () => void;
    errorMessage: string;
}

/**
 * A form component for creating or editing a task.
 * @function TaskForm
 * @param props - The props for the component.
 * @returns {JSX.Element} - The rendered task form.
 */
const TaskForm: React.FC<TaskFormProps> = ({ task, onSave, onCancel, errorMessage }): JSX.Element => {
    const dispatch = useAppDispatch();
    const [taskName, setTaskName] = useState(task.name);
    const [taskDescription, setTaskDescription] = useState(task.description);
    const [taskNameError, setTaskNameError] = useState("");

    const handleSave = async () => {
        if (taskName.trim() === "") {
            setTaskNameError("Task name cannot be empty");
            logger.warn("Attempted to save task with empty name");
            dispatch(setFeedback({ message: "Attempted to save task with empty name", severity: "error" }));
            return;
        }
        setTaskNameError("");

        try {
            await onSave({ ...task, name: taskName, description: taskDescription });
            logger.info("Task saved successfully");
            dispatch(setFeedback({ message: "Task saved successfully", severity: "success" }));
        } catch (error) {
            logger.error("Failed to save task", error);
            dispatch(setFeedback({ message: "Failed to save task", severity: "error" }));
        }
    };

    return (
        <Box component="form">
            <TextField
                label="Task Name"
                value={ taskName }
                onChange={ (e) => setTaskName(e.target.value) }
                variant="outlined"
                error={ !!taskNameError }
                helperText={ taskNameError }
            />
            <TextareaAutosize
                minRows={ 4 }
                placeholder="Task Description"
                value={ taskDescription }
                onChange={ (e) => setTaskDescription(e.target.value) }
                style={ { width: "100%", padding: "8px", fontSize: "16px" } }
            />
            { errorMessage && (
                <Typography color="error" variant="body2">
                    { errorMessage }
                </Typography>
            ) }
            <Box sx={ { display: "flex", justifyContent: "flex-end", gap: 1 } }>
                <Button variant="outlined" color="secondary" onClick={ onCancel }>Cancel</Button>
                <Button variant="contained" color="primary" onClick={ handleSave }>Save</Button>
            </Box>
        </Box>
    );
};

export default TaskForm;
