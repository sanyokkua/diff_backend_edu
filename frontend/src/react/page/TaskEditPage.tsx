import { Box, Typography }         from "@mui/material";
import { FC, useEffect, useState } from "react";
import { useNavigate }             from "react-router-dom";
import {
    parseErrorMessage,
    setAppBarHeader,
    TaskDto,
    updateTask,
    UpdateTaskRequest,
    useAppDispatch,
    useAppSelector
}                                  from "../../core";
import TaskForm                    from "../component/TaskForm";


const TaskEditPage: FC = () => {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();
    const [errorMsg, setErrorMsg] = useState("");
    const { userId, taskId, taskName, taskDescription } = useAppSelector((state) => state.globals);
    const task: TaskDto = { taskId: taskId, name: taskName, description: taskDescription, userId: userId };

    useEffect(() => {
        dispatch(setAppBarHeader("Edit Task"));
    }, [dispatch]);

    const onSubmit = async (data: TaskDto) => {
        const updateReq: UpdateTaskRequest = {
            userId: userId,
            taskId: taskId,
            taskUpdateDTO: {
                name: data.name,
                description: data.description
            }
        };
        try {
            await dispatch(updateTask(updateReq)).unwrap();
            navigate("/dashboard");
        } catch (e) {
            const errMsg = parseErrorMessage(e);
            setErrorMsg(errMsg);
        }
    };

    const onCancel = () => {
        navigate("/dashboard");
    };

    return <Box
        sx={ {
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
            width: "80%",
            margin: "0 auto",
            padding: 2
        } }
    >
        <Typography variant="h4" gutterBottom>
            Update task
        </Typography>

        <TaskForm task={ task } onSave={ onSubmit } onCancel={ onCancel } error={ errorMsg }/>
    </Box>;
};

export default TaskEditPage;
