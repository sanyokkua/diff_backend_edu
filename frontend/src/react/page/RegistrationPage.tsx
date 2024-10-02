import { FC, useEffect }                                                  from "react";
import { useNavigate }                                                    from "react-router-dom";
import { registerUser, setAppBarHeader, useAppDispatch, UserCreationDTO } from "../../core";
import FormGeneric, { FormResponse }                                      from "../component/FormGeneric";


const RegistrationPage: FC = () => {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

    useEffect(() => {
        dispatch(setAppBarHeader("Registration"));
    }, [dispatch]);

    const onSubmit = async (data: FormResponse) => {
        const regReq: UserCreationDTO = {
            email: data.email,
            password: data.password,
            passwordConfirmation: data.confirmPassword ?? ""
        };
        await dispatch(registerUser(regReq)).unwrap();
        navigate("/dashboard");
    };

    return <FormGeneric formType={ "Register" } onSubmit={ onSubmit }/>;
};

export default RegistrationPage;
