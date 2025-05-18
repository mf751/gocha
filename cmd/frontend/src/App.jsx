import { Routes, Route, useLocation, useNavigate } from "react-router-dom";
import Login from "./pages/login/login.jsx";
import Signup from "./pages/login/singup.jsx";
import Home from "./pages/home/index.jsx";
import Theme from "./pages/theme/index.jsx";
import Profile from "./pages/profile/index.jsx";
import Chat from "./pages/chat/index.jsx";
import RequireAuth from "./helpers/middleware.jsx";
import { useDispatch, useSelector } from "react-redux";
import APIURL from "./api.js";
import { setLoggedIn, setUser } from "./store/slices/user.js";
import { setChats, setLoaded } from "./store/slices/chats.js";
import { useEffect, useState } from "react";
import Nav from "./pages/nav/index.jsx";
import { useRef } from "react";

function App() {
  const location = useLocation();
  const dispatch = useDispatch();
  const loggedIn = useSelector((state) => state.user.loggedIn);
  const user = useSelector((state) => state.user.user);
  const chats = useSelector((state) => state.chats.chats);
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const wsRef = useRef(null);

  // runs on every reroute
  useEffect(() => {
    const expiry = localStorage.getItem("expiry");
    if (expiry === null || new Date(expiry) < new Date()) {
      localStorage.setItem("expiry", "");
      localStorage.setItem("authToken", "");
      navigate("/login", { replace: true });
      setLoading(false);
    }
    if (!loggedIn && localStorage.getItem("expiry") !== "") {
      const token = localStorage.getItem("authToken");
      (async () => {
        try {
          const res = await fetch(`${APIURL}/v1/user`, {
            headers: { Authorization: `Bearer ${token}` },
          });
          const data = await res.json();
          if (res.status !== 200) {
            navigate("/login", { replace: true });
            setLoading(false);
          }

          dispatch(setUser(data.user));
          dispatch(setLoggedIn(true));
          navigate(window.location.pathname, { replace: true });
          setLoading(false);
        } catch (error) {
          setLoading(false);
        }
      })();
    }
  }, [location]);

  // Only runs once per user
  useEffect(() => {
    if (Object.keys(user).length !== 0) {
      let socket;
      (async () => {
        try {
          const token = localStorage.getItem("authToken");
          const res = await fetch(`${APIURL}/v1/chats`, {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          });
          const data = await res.json();
          dispatch(setChats(data.data));
          dispatch(setLoaded(true));
          // only initiate the ws and the next useEffect will set the onmessage on every chats change
          socket = new WebSocket(`${APIURL}/v1/ws?token=${token}`);
          wsRef.current = socket;
        } catch (error) {
          console.log(error);
        }
        return () => {
          if (socket) socket.close();
        };
      })();
    }
  }, [user]);

  // runs on every chats value change so that it resets the onmessage function correctly!
  useEffect(() => {
    if (!wsRef.current) return;
    wsRef.current.onmessage = (evt) => {
      const wsData = JSON.parse(evt.data);
      const newChats = chats.map((obj) => {
        if (obj.chat.id != wsData.payload.chat_id) return obj;
        return {
          chat: obj.chat,
          members: obj.members,
          last_message: {
            message: {
              chat_id: wsData.payload.chat_id,
              content: wsData.payload.message,
              id: wsData.payload.id,
              sent: wsData.payload.sent,
              user_name: wsData.payload.user_name,
              type:
                wsData.type === "new_message"
                  ? 1
                  : wsData.type === "joined_message"
                    ? 50
                    : 51,
            },
            user: {
              id: wsData.payload.from,
              name: wsData.payload.user_name,
            },
          },
        };
      });
      dispatch(setChats(newChats));
    };
  }, [chats]);

  if (loading) {
    return (
      <div
        style={{
          height: "100dvh",
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
        }}
      >
        <h1>Loading</h1>
      </div>
    );
  }
  return (
    <div className="parent">
      <Routes>
        <Route
          path="/login"
          element={
            <>
              <Theme />
              <Login />
            </>
          }
        />
        <Route
          path="/signup"
          element={
            <>
              <Theme />
              <Signup />
            </>
          }
        />
        <Route
          path="/"
          element={
            <RequireAuth>
              <Theme />
              <Nav />
              <Home />
            </RequireAuth>
          }
        />
        <Route
          path="/profile"
          element={
            <RequireAuth>
              <Theme />
              <Nav />
              <Profile />
            </RequireAuth>
          }
        />
        <Route
          path="/chat/:chatID"
          element={
            <RequireAuth>
              <Chat />
            </RequireAuth>
          }
        />
      </Routes>
    </div>
  );
}

export default App;
