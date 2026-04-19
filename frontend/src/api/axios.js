// frontend/src/api/axios.js
import axios from "axios";
import { useUserStore } from "../store/user";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '/api', // fallback для nginx proxy
});

// Интерцептор запроса: токен читается ДИНАМИЧЕСКИ при каждом запросе
api.interceptors.request.use((config) => {
  const token = useUserStore.getState().token; // ← читаем свежее значение
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
}, (error) => {
  return Promise.reject(error);
});

export default api;