import axios from "axios";

const api = axios.create({
  baseURL: `${process.env.NEXT_PUBLIC_BASE_URL || process.env.REACT_APP_BASE_URL}/api`,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

// login api
export const loginAPI = (email, password) =>
  api.post("/auth/login", { email, password });
// get all  folders
export const getFolders = () => api.get("/folders");
// get all photos
export const getPhotos = () => api.get("/photos");
