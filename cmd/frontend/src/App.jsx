import { Routes, Route } from "react-router-dom";
import Login from "./pages/login/login.jsx";
import Signup from "./pages/login/singup.jsx";
import Home from "./pages/home/index.jsx";
import Theme from "./pages/theme/index.jsx";

function App() {
  return (
    <div className="parent">
      <Theme />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<Signup />} />
      </Routes>
    </div>
  );
}

export default App;
