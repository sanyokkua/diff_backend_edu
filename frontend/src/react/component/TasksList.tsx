import DeleteIcon                                                                      from "@mui/icons-material/Delete";
import { Divider, IconButton, List, ListItem, ListItemButton, ListItemText, Skeleton } from "@mui/material";
import { FC, MouseEvent, useCallback }                                                 from "react";
import { TaskDto }                                                                     from "../../core";


interface TasksListProps {
    tasks: TaskDto[];
    onTaskSelected: (dto: TaskDto) => void;
    onTaskDeleted: (dto: TaskDto) => void;
}

const TasksList: FC<TasksListProps> = ({ tasks, onTaskSelected, onTaskDeleted }) => {
    const handleOnTaskClicked = useCallback(
        (task: TaskDto) => (event: MouseEvent) => {
            event.stopPropagation();
            onTaskSelected(task);
        }, [onTaskSelected]);

    const handleOnTaskDeleteClicked = useCallback(
        (task: TaskDto) => (event: MouseEvent) => {
            event.stopPropagation();
            onTaskDeleted(task);
        }, [onTaskDeleted]);

    const mapTaskToListItem = (taskDto: TaskDto, index: number) => (
        <div key={ taskDto.taskId }>
            <ListItem disablePadding onClick={ handleOnTaskClicked(taskDto) }>
                <ListItemButton>
                    <ListItemText primary={ taskDto.name }/>
                </ListItemButton>
                <IconButton edge="end" aria-label="delete" onClick={ handleOnTaskDeleteClicked(taskDto) }
                            sx={ { marginRight: 3 } }>
                    <DeleteIcon sx={ { color: "red" } }/>
                </IconButton>
            </ListItem>
            { index < tasks.length - 1 && <Divider/> }
        </div>
    );

    const skeleton = <Skeleton variant="rectangular" height={ 200 }/>;
    const content = <List>{ tasks.map(mapTaskToListItem) }</List>;

    const toDisplay = tasks.length > 0 ? content : skeleton;
    return <>{ toDisplay }</>;
};

export default TasksList;
