import { useSelector } from "react-redux";
import { Navigate } from "react-router-dom";

export default function RequireAuth({ children }) {
  const isAuthenticated = useSelector((state) => state.user.loggedIn);
  console.log("called");
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return children;
}
