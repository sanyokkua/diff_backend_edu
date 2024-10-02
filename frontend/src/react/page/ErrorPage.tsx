import { FC, useEffect }                   from "react";
import { useRouteError }                   from "react-router-dom";
import { setAppBarHeader, useAppDispatch } from "../../core";


const ErrorPage: FC = () => {
    const dispatch = useAppDispatch();

    useEffect(() => {
        dispatch(setAppBarHeader("Error"));
    }, [dispatch]);

    const error = useRouteError();
    const errorAsObj: { statusText?: string, message?: string } = error as { statusText?: string, message?: string };

    return <>
        <h1>Oops!</h1>
        <p>Sorry, an unexpected error has occurred.</p>
        <p>
            <i>{ errorAsObj?.statusText ?? errorAsObj?.message }</i>
        </p>
    </>;
};

export default ErrorPage;