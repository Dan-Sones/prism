import axios from "axios";

declare global {
  interface Window {
    __APP_CONFIG__?: { apiBaseUrl?: string };
  }
}

const baseURL =
  window.__APP_CONFIG__?.apiBaseUrl ?? "http://localhost:8081/api";

export const axiosClient = axios.create({
  baseURL,
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 10_000,
});
