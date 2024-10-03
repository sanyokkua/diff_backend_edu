import { Box, Typography }              from "@mui/material";
import { FC, JSX, useEffect, useState } from "react";
import { useNavigate }                  from "react-router-dom";
import {
    LogLevel,
    parseErrorMessage,
    setHeaderTitle,
    TaskDto,
    updateTask,
    UpdateTaskPayload,
    useAppDispatch,
    useAppSelector
}                                       from "../../core";
import { TaskForm }                     from "../component";


const logger = LogLevel.getLogger("TaskEditPage");

/**
 * TaskEditPage component.
 *
 * This component is responsible for rendering the task edit page and handling the update of existing tasks.
 * It sets the header title to "Edit Task" and manages the form submission process.
 *
 * @component
 * @returns {JSX.Element} The rendered task edit page component.
 */
const TaskEditPage: FC = (): JSX.Element => {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();
    const [errorMsg, setErrorMsg] = useState<string>("");
    const { userId } = useAppSelector((state) => state.users);
    const { taskId, taskName, taskDescription } = useAppSelector((state) => state.tasks);
    const task: TaskDto = { taskId: taskId, name: taskName, description: taskDescription, userId: userId };

    useEffect(() => {
        // Set the header title to "Edit Task" when the component mounts.
        dispatch(setHeaderTitle("Edit Task"));
    }, [dispatch]);

    /**
     * Handles the form submission.
     *
     * @param {TaskDto} data - The form data containing task details.
     */
    const onSubmit = async (data: TaskDto) => {
        logger.debug("Submit", data);
        const updateReq: UpdateTaskPayload = {
            userId: userId,
            taskId: taskId,
            taskData: {
                name: data.name,
                description: data.description
            }
        };
        try {
            await dispatch(updateTask(updateReq)).unwrap();
            navigate("/dashboard");
        } catch (e) {
            const errMsg = parseErrorMessage(e, "Failed to update task");
            setErrorMsg(errMsg);
        }
    };

    /**
     * Handles the cancel action.
     */
    const onCancel = () => {
        navigate("/dashboard");
    };

    return (
        <Box>
            <Typography variant="h4" gutterBottom>
                Update task
            </Typography>
            <TaskForm task={ task } onSave={ onSubmit } onCancel={ onCancel } errorMessage={ errorMsg }/>
        </Box>
    );
};

export default TaskEditPage;
