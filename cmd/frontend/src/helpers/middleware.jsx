import { useSelector } from "react-redux";
import { Navigate } from "react-router-dom";

export default function RequireAuth({ children }) {
  const isAuthenticated = useSelector((state) => state.user.loggedIn);
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return children;
}
