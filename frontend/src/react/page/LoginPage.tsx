import { FC, useEffect }                                            from "react";
import { useNavigate }                                              from "react-router-dom";
import { loginUser, setAppBarHeader, useAppDispatch, UserLoginDto } from "../../core";
import FormGeneric, { FormResponse }                                from "../component/FormGeneric";


const LoginPage: FC = () => {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

    useEffect(() => {
        dispatch(setAppBarHeader("Login"));
    }, [dispatch]);

    const onSubmit = async (data: FormResponse) => {
        const loginReq: UserLoginDto = {
            email: data.email,
            password: data.password
        };
        await dispatch(loginUser(loginReq)).unwrap();
        navigate("/dashboard");
    };

    return <FormGeneric formType={ "Login" } onSubmit={ onSubmit }/>;
};

export default LoginPage;
