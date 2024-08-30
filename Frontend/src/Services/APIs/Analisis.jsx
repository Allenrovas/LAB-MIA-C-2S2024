import axios from 'axios';

const instance = axios.create({
    baseURL: 'http://localhost:3000/',
    timeout: 15000,
    headers: {
        'Content-Type': 'application/json',
    },
});

export const analisis = async (peticion) => {
    try {
        const { data } = await instance.post('command/', { peticion });
        return data;
    } catch (error) {
        console.error('Error during API call:', error);
        throw error;
    }
};
