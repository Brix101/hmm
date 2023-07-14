import { BASE_URL } from "@/constant/server.constant";
import axios from "axios";

const apiClient = axios.create({
  baseURL: BASE_URL,
  withCredentials: true,
});

export { apiClient };
