import { Box, Button, CircularProgress } from "@mui/material";
import { FC, JSX, useEffect, useState }  from "react";
import { useNavigate }                   from "react-router-dom";
import {
    deleteTask,
    DeleteTaskPayload,
    fetchTasks,
    LogLevel,
    setFeedback,
    setHeaderTitle,
    setTaskDescription,
    setTaskId,
    setTaskName,
    TaskDto,
    useAppDispatch,
    useAppSelector
}                                        from "../../core";
import { ConfirmationDialog, TasksList } from "../component";


const logger = LogLevel.getLogger("Dashboard");

/**
 * Dashboard component.
 *
 * This component is responsible for displaying the user's tasks and providing functionality
 * to add, edit, and delete tasks. It uses Redux for state management and logs various actions.
 *
 * @component
 * @returns {JSX.Element} The rendered dashboard component.
 */
const Dashboard: FC = (): JSX.Element => {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

    const [isDialogOpen, setIsDialogOpen] = useState<boolean>(false);
    const [taskToDelete, setTaskToDelete] = useState<DeleteTaskPayload | null>(null);

    const { tasksList, isTaskLoading, taskError } = useAppSelector((state) => state.tasks);
    const { userId } = useAppSelector((state) => state.users);

    useEffect(() => {
        // Set the header title to "Dashboard" when the component mounts.
        dispatch(setHeaderTitle("Dashboard"));
        if (userId) {
            // Fetch tasks for the user if userId is available.
            dispatch(fetchTasks(userId));
        }
    }, [dispatch, userId]);

    useEffect(() => {
        if (taskError) {
            // Display feedback and log an error if there is a task error.
            dispatch(setFeedback({ message: taskError, severity: "error" }));
            logger.error(`Task error: ${ taskError }`);
        }
    }, [dispatch, taskError]);

    /**
     * Handles the click event on a task.
     *
     * @param {TaskDto} task - The task object.
     */
    const handleTaskClick = (task: TaskDto) => {
        dispatch(setTaskId(task.taskId));
        dispatch(setTaskName(task.name));
        dispatch(setTaskDescription(task.description));
        navigate("/dashboard/edit");
    };

    /**
     * Handles the click event to delete a task.
     *
     * @param {TaskDto} task - The task object.
     */
    const handleDeleteClick = (task: TaskDto) => {
        setTaskToDelete({ userId, taskId: task.taskId });
        setIsDialogOpen(true);
    };

    /**
     * Handles the deletion of a task.
     */
    const handleDelete = async () => {
        if (taskToDelete) {
            try {
                await dispatch(deleteTask(taskToDelete));
                await dispatch(fetchTasks(userId));
                logger.info(`Task deleted: ${ taskToDelete.taskId }`);
            } catch (error) {
                logger.error(`Failed to delete task: ${ error }`);
                dispatch(setFeedback({ message: "Failed to delete task", severity: "error" }));
            }
        }
        resetDeleteState();
    };

    /**
     * Resets the state related to task deletion.
     */
    const resetDeleteState = () => {
        setIsDialogOpen(false);
        setTaskToDelete(null);
    };

    /**
     * Handles the click event to add a new task.
     */
    const handleAddTaskClick = () => {
        navigate("/dashboard/new");
    };

    logger.debug("Dashboard render");
    return (
        <>
            <ConfirmationDialog
                title="Delete Task"
                message="Are you sure you want to delete this task? This action cannot be undone."
                isOpen={ isDialogOpen }
                onConfirm={ handleDelete }
                onCancel={ resetDeleteState }
                onClose={ resetDeleteState }
            />
            { isTaskLoading ? (
                <CircularProgress/>
            ) : (
                  <Box display="flex" justifyContent="center" alignItems="center" width="80%" mx="auto" my={ 2 }>
                      <Box width="100%">
                          <Box display="flex" justifyContent="center" my={ 2 }>
                              <Button variant="contained" color="success" onClick={ handleAddTaskClick }>
                                  Add Task
                              </Button>
                          </Box>
                          <TasksList tasks={ tasksList } onTaskSelect={ handleTaskClick }
                                     onTaskDelete={ handleDeleteClick }/>
                      </Box>
                  </Box>
              ) }
        </>
    );
};

export default Dashboard;
