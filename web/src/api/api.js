import axios from "axios";

export const apiClient = axios.create({
  baseURL: `http://${import.meta.env.VITE_API_URL}`,
  headers: {
    "Content-Type": "application/json",
  },
});
