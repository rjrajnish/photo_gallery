import Login from "@/components/Login";
import Dashboard from "@/components/Dashboard";
import { AuthProvider, useAuth } from "./authProvider";

function Main() {
  const { token } = useAuth();
  return token ? <Dashboard /> : <Login />;
}

export default function Home() {
  return (
    <AuthProvider>
      <Main />
    </AuthProvider>
  );
}
