import { Routes, Route, useLocation, useNavigate } from "react-router-dom";
import Login from "./pages/login/login.jsx";
import Signup from "./pages/login/singup.jsx";
import Home from "./pages/home/index.jsx";
import Theme from "./pages/theme/index.jsx";
import Profile from "./pages/profile/index.jsx";
import RequireAuth from "./helpers/middleware.jsx";
import { useDispatch, useSelector } from "react-redux";
import APIURL from "./api.js";
import { setLoggedIn, setUser } from "./store/slices/user.js";
import { useEffect, useState } from "react";

function App() {
  const location = useLocation();
  const dispatch = useDispatch();
  const loggedIn = useSelector((state) => state.user.loggedIn);
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const expiry = localStorage.getItem("expiry");
    if (expiry === null || new Date(expiry) < new Date())
      navigate("/login", { replace: true });
    if (!loggedIn) {
      const token = localStorage.getItem("authToken");
      (async () => {
        try {
          const res = await fetch(`${APIURL}/v1/user`, {
            headers: { Authorization: `Bearer ${token}` },
          });
          const data = await res.json();
          dispatch(setUser(data.user));
          dispatch(setLoggedIn(true));
          navigate(window.location.pathname, { replace: true });
          setLoading(false);
        } catch (error) {
          console.log(error);
          setLoading(false);
          navigate("/login", { replace: true });
        }
      })();
    }
  }, [location, dispatch, loggedIn]);
  if (loading) {
    return <h1>Loading</h1>;
  }
  return (
    <div className="parent">
      <Theme />
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<Signup />} />
        <Route
          path="/"
          element={
            <RequireAuth>
              <Home />
            </RequireAuth>
          }
        />
        <Route
          path="/profile"
          element={
            <RequireAuth>
              <Profile />
            </RequireAuth>
          }
        />
      </Routes>
    </div>
  );
}

export default App;
