import DeleteIcon                                                                      from "@mui/icons-material/Delete";
import { Divider, IconButton, List, ListItem, ListItemButton, ListItemText, Skeleton } from "@mui/material";
import { FC, JSX, MouseEvent, useCallback }                                            from "react";
import { LogLevel, setFeedback, TaskDto, useAppDispatch }                              from "../../core";


const logger = LogLevel.getLogger("TasksList");

/**
 * Props for the TasksList component.
 * @property {TaskDto[]} tasks - The list of tasks to display.
 * @property {Function} onTaskSelect - Function to call when a task is selected.
 * @property {Function} onTaskDelete - Function to call when a task is deleted.
 */
interface TasksListProps {
    tasks: TaskDto[];
    onTaskSelect: (task: TaskDto) => void;
    onTaskDelete: (task: TaskDto) => void;
}

/**
 * A component to display a list of tasks with options to select and delete each task.
 * @function TasksList
 * @param props - The props for the component.
 * @returns {JSX.Element} - The rendered tasks list.
 */
const TasksList: FC<TasksListProps> = ({ tasks, onTaskSelect, onTaskDelete }): JSX.Element => {
    const dispatch = useAppDispatch();

    /**
     * Handles the click event for selecting a task.
     * @function handleTaskClick
     * @param {TaskDto} task - The task to select.
     * @returns {Function} - The event handler function.
     */
    const handleTaskClick = useCallback(
        (task: TaskDto) => (event: MouseEvent) => {
            event.stopPropagation();
            logger.debug("Task clicked:", task);
            try {
                onTaskSelect(task);
            } catch (error) {
                logger.error("Error selecting task:", error);
                dispatch(setFeedback({ message: "Error selecting task", severity: "error" }));
            }
        },
        [dispatch, onTaskSelect]
    );

    /**
     * Handles the click event for deleting a task.
     * @function handleTaskDeleteClick
     * @param {TaskDto} task - The task to delete.
     * @returns {Function} - The event handler function.
     */
    const handleTaskDeleteClick = useCallback(
        (task: TaskDto) => (event: MouseEvent) => {
            event.stopPropagation();
            logger.debug("Task delete clicked:", task);
            try {
                onTaskDelete(task);
            } catch (error) {
                logger.error("Error deleting task:", error);
                dispatch(setFeedback({ message: "Error deleting task", severity: "error" }));
            }
        },
        [dispatch, onTaskDelete]
    );

    /**
     * Renders a task item.
     * @function renderTaskItem
     * @param {TaskDto} task - The task to render.
     * @param {number} index - The index of the task in the list.
     * @returns {JSX.Element} - The rendered task item.
     */
    const renderTaskItem = (task: TaskDto, index: number): JSX.Element => (
        <div key={ task.taskId }>
            <ListItem disablePadding onClick={ handleTaskClick(task) }>
                <ListItemButton>
                    <ListItemText primary={ task.name }/>
                </ListItemButton>
                <IconButton edge="end" aria-label="delete" onClick={ handleTaskDeleteClick(task) }
                            sx={ { marginRight: 3 } }>
                    <DeleteIcon sx={ { color: "red" } }/>
                </IconButton>
            </ListItem>
            { index < tasks.length - 1 && <Divider/> }
        </div>
    );

    const skeleton = <Skeleton variant="rectangular" height={ 200 }/>;
    const items = <List>{ tasks.map(renderTaskItem) }</List>;
    const content = tasks.length > 0 ? items : skeleton;

    return <>{ content }</>;
};

export default TasksList;
