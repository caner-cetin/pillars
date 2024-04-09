import axios from "axios";

const restURL = 'http://localhost:1323/api/'
const wsURL = 'http://0.0.0.0:1323/ws/'
const httpClient = axios.create({
    baseURL: restURL,
    timeout: 1000,
});


export { httpClient };
export { restURL };
export { wsURL };