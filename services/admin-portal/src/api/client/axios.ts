import axios from "axios";

export const axiosClient = axios.create({
  baseURL: "http://localhost:8081/api",
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 10_000,
});
