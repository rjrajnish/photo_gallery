import { createContext, useContext, useState, useEffect } from "react";
import { loginAPI } from "./api/api";

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null); // store user info
  const [token, setToken] = useState(null); // JWT token

  // 🔹 Restore user/token from localStorage on refresh
  useEffect(() => {
    const savedToken = localStorage.getItem("token");

    if (savedToken) {
      setToken(savedToken);
    }
  }, [token]);

  // 🔹 Login using API
  const login = async (email, password) => {
    try {
      const res = await loginAPI(email, password);

      if (res?.data?.token) {
        setToken(res?.data.token);

        // save to localStorage
        localStorage.setItem("token", res.data.token);

        return true;
      }
      return false;
    } catch (err) {
      console.error("Login failed:", err);
      return false;
    }
  };

  // 🔹 Logout
  const logout = () => {
    setToken(null);
    localStorage.removeItem("token");
  };

  return (
    <AuthContext.Provider value={{ token, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
