import { Box, Button, CircularProgress } from "@mui/material";
import { FC, useEffect, useState }       from "react";
import { useNavigate }                   from "react-router-dom";
import {
    deleteTask,
    DeleteUserTaskRequest,
    getTasks,
    setAppBarHeader,
    setTaskDescription,
    setTaskId,
    setTaskName,
    TaskDto,
    useAppDispatch,
    useAppSelector
}                                        from "../../core";
import ConfirmationDialog                from "../component/ConfirmationDialog";
import FeedbackSnackbar                  from "../component/FeedbackSnackbar";
import TasksList                         from "../component/TasksList";


const Dashboard: FC = () => {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

    const [openDialog, setOpenDialog] = useState<boolean>(false);
    const [deleteTaskDto, setDeleteTaskDto] = useState<DeleteUserTaskRequest | null>(null);

    const { userId, tasksList, appIsLoading, appError } = useAppSelector((state) => state.globals);

    useEffect(() => {
        dispatch(setAppBarHeader("Dashboard"));
        if (userId) {
            dispatch(getTasks(userId)); // Fetch tasks on mount
        }
    }, [dispatch, userId]);

    const handleTaskClick = (task: TaskDto) => {
        dispatch(setTaskId(task.taskId));
        dispatch(setTaskName(task.name));
        dispatch(setTaskDescription(task.description));
        navigate("/dashboard/edit");
    };

    const handleDeleteClick = (task: TaskDto) => {
        const delReq: DeleteUserTaskRequest = {
            userId: userId,
            taskId: task.taskId
        };
        setDeleteTaskDto(delReq);
        setOpenDialog(true);
    };

    const handleDelete = async () => {
        if (deleteTaskDto) {
            await dispatch(deleteTask(deleteTaskDto));
            await dispatch(getTasks(userId)); // Re-Fetch tasks after delete
        }
        resetDelete();
    };

    const resetDelete = () => {
        setOpenDialog(false);
        setDeleteTaskDto(null);
    };

    const handleAddTaskClick = () => {
        navigate("/dashboard/new");
    };

    return (
        <>
            <ConfirmationDialog
                title="Delete Task"
                content="Are you sure you want to delete this task? This action cannot be undone."
                open={ openDialog }
                onConfirm={ handleDelete }
                onCancel={ resetDelete }
                onClose={ resetDelete }
            />
            { appIsLoading && <CircularProgress/> }
            { !appIsLoading && (
                <Box display="flex" justifyContent="center" alignItems="center" width="80%" mx="auto" my={ 2 }>
                    <Box width="100%">
                        <Box display="flex" justifyContent="center" my={ 2 }>
                            <Button variant="contained" color="success" onClick={ handleAddTaskClick }>
                                Add Task
                            </Button>
                        </Box>
                        <TasksList tasks={ tasksList } onTaskSelected={ handleTaskClick }
                                   onTaskDeleted={ handleDeleteClick }/>
                    </Box>
                </Box>
            ) }
            { appError && (
                <FeedbackSnackbar
                    feedbackMessage={ { message: appError, severity: "error" } }
                    onClose={ () => {
                    } }
                />
            ) }
        </>
    );
};

export default Dashboard;
