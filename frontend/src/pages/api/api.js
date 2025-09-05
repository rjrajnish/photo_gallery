import axios from "axios";

const api = axios.create({
  baseURL: `${
    process.env.NEXT_PUBLIC_BASE_URL || process.env.REACT_APP_BASE_URL
  }/api`,
});

// ✅ Request Interceptor (add token)
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) config.headers.Authorization = `Bearer ${token}`;
    return config;
  },
  (error) => {
    // Request error
    return Promise.reject(error);
  }
);

// ✅ Response Interceptor (handle errors globally)
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      // Server responded with a status other than 2xx
      console.error("API Error:", error.response.data);
      return Promise.reject(error.response.data);
    } else if (error.request) {
      // Request was made but no response
      console.error("Network Error:", error.request);
      return Promise.reject({ message: "Network error, please try again." });
    } else {
      // Something happened in setting up the request
      console.error("Unexpected Error:", error.message);
      return Promise.reject({ message: "Unexpected error occurred." });
    }
  }
);

// ✅ Example API function
export const loginAPI = async (email, password) => {
  try {
    const res = await api.post("/auth/login", { email, password });
    return res
  } catch (err) {
    throw err; // re-throw so component can catch
  }
};
// create new user
export const createUser = async (name, email, password) => {
  try {
    const res = await api.post("/auth/register", { name, email, password });
    return res
  } catch (error) {
    throw error;
  }
}
export const getFolders = async () => {
  try {
    const res = await api.get("/folders");
    return res
  } catch (error) {
    throw error;
  }
};

// get all photos

export const getPhotos = async () => {
  try {
    const res = await api.get("/photos");
    return res
  } catch (error) {
    throw error;
  }
};

// create folder

export const createFolder = async (name) => {
  try {
    const res = await api.post("/folders", { name });
    return res
  } catch (error) {
    throw error;
  }
};
