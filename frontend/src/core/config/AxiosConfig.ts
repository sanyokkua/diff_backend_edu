import axios from "axios";


const baseUrl = "http://localhost:8080";

const axiosClient = axios.create(
    {
        baseURL: `${ baseUrl }/api/v1`,
        headers: {
            "Content-Type": "application/json"
        },
        withCredentials: true
    }
);

export default axiosClient;