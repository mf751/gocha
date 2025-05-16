import { useDispatch, useSelector } from "react-redux";
import "./styles.css";
import { Logout } from "./../../helpers/auth.js";

export default function Profile() {
  const user = useSelector((state) => state.user.user);
  const dispatch = useDispatch();
  return (
    <div className="profile">
      <h1>Profile</h1>
      <div className="info">
        <div className="card tag">Name</div>
        <div className="card">{user.name}</div>
        <div className="card tag">Email</div>
        <div className="card">{user.email}</div>
        <div className="card tag">Created</div>
        <div className="card">
          {new Date(user.created_at).toLocaleDateString()}
        </div>
        <div className="card tag">Activated</div>
        <div className="card">{user.activated ? "Yes" : "No"}</div>
        <div className="card tag">ID</div>
        <div className="card id">{user.id}</div>
      </div>
      <button className="logout" onClick={() => Logout(dispatch)}>
        Logout
      </button>
    </div>
  );
}
