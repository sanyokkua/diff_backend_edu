import { Box, Typography }         from "@mui/material";
import { FC, useEffect, useState } from "react";
import { useNavigate }             from "react-router-dom";
import {
    createTask,
    CreateTaskRequest,
    parseErrorMessage,
    setAppBarHeader,
    TaskDto,
    useAppDispatch,
    useAppSelector
}                                  from "../../core";
import TaskForm                    from "../component/TaskForm";


const TaskCreatePage: FC = () => {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();
    const [errorMsg, setErrorMsg] = useState<string>("");
    const { userId } = useAppSelector((state) => state.globals);
    const task: TaskDto = { taskId: 0, name: "", description: "", userId: userId };

    useEffect(() => {
        dispatch(setAppBarHeader("Create Task"));
    }, [dispatch]);

    const onSubmit = async (data: TaskDto) => {
        const updateReq: CreateTaskRequest = {
            userId: userId,
            taskCreationDTO: {
                name: data.name,
                description: data.description
            }
        };
        try {
            await dispatch(createTask(updateReq)).unwrap();
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
            Create new task
        </Typography>

        <TaskForm task={ task } onSave={ onSubmit } onCancel={ onCancel } error={ errorMsg }/>
    </Box>;
};

export default TaskCreatePage;
