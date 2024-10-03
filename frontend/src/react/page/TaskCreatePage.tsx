import { Box, Typography }              from "@mui/material";
import { FC, JSX, useEffect, useState } from "react";
import { useNavigate }                  from "react-router-dom";
import {
    createTask,
    CreateTaskPayload,
    LogLevel,
    parseErrorMessage,
    setHeaderTitle,
    TaskDto,
    useAppDispatch,
    useAppSelector
}                                       from "../../core";
import { TaskForm }                     from "../component";


const logger = LogLevel.getLogger("TaskCreatePage");

/**
 * TaskCreatePage component.
 *
 * This component is responsible for rendering the task creation page and handling the creation of new tasks.
 * It sets the header title to "Create Task" and manages the form submission process.
 *
 * @component
 * @returns {JSX.Element} The rendered task creation page component.
 */
const TaskCreatePage: FC = (): JSX.Element => {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();
    const [errorMsg, setErrorMsg] = useState<string>("");
    const { userId } = useAppSelector((state) => state.users);
    const task: TaskDto = { taskId: 0, name: "", description: "", userId: userId };

    useEffect(() => {
        // Set the header title to "Create Task" when the component mounts.
        dispatch(setHeaderTitle("Create Task"));
    }, [dispatch]);

    /**
     * Handles the form submission.
     *
     * @param {TaskDto} data - The form data containing task details.
     */
    const onSubmit = async (data: TaskDto) => {
        logger.debug("Submit", data);
        const updateReq: CreateTaskPayload = {
            userId: userId,
            taskData: {
                name: data.name,
                description: data.description
            }
        };
        try {
            await dispatch(createTask(updateReq)).unwrap();
            navigate("/dashboard");
        } catch (e) {
            const errMsg = parseErrorMessage(e, "Failed to create task");
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
                Create new task
            </Typography>
            <TaskForm task={ task } onSave={ onSubmit } onCancel={ onCancel } errorMessage={ errorMsg }/>
        </Box>
    );
};

export default TaskCreatePage;
