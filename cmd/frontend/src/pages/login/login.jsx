import { Link, useNavigate } from "react-router-dom";
import "./styles.css";
import { useState, useRef } from "react";
import APIURL from "./../../api.js";
import { SetAuthInfo } from "../../helpers/auth.js";
import { useDispatch } from "react-redux";

export default function Login() {
  const [form, setForm] = useState({ email: "", password: "" });
  const [failed, setFailed] = useState({});
  const btnRef = useRef();
  const navigate = useNavigate();
  const dispatch = useDispatch();

  function login(event) {
    const email = event.get("email");
    const password = event.get("password");
    setForm({ email: email, password: password });
    fetch(`${APIURL}/v1/tokens/authentication`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        email: email,
        password: password,
      }),
    })
      .then((res) => res.json())
      .then((res) => {
        if (res["error"] !== undefined) {
          let key = Object.keys(res["error"])[0];
          if (key === "0") {
            setFailed({ type: "general", msg: res["error"] });
            return;
          }
          setFailed({ type: key, msg: res["error"][key] });
          return;
        }
        if (res["authentication_token"] === undefined) {
          setFailed({ type: "general", msg: "something went wrong" });
          return;
        }
        setFailed("");
        btnRef.current.style.backgroundColor = "#41dc83";
        btnRef.current.disabled = true;
        btnRef.current.textContent = "Logged in Succuessfully!";
        SetAuthInfo(res, dispatch);
        setTimeout(() => navigate("/profile"), 1000);
      });
  }
  return (
    <div className="login">
      <div className="page-title">
        <h1>Login</h1>
      </div>
      <form action={login} className="login-form">
        <div className="inputs">
          <input
            id="email"
            type="email"
            placeholder="Email"
            name="email"
            required
            className={failed.type === "email" ? "failed" : ""}
            defaultValue={form.email}
          />
          {failed.type === "emali" && (
            <span className="failed-msg">{failed.msg}</span>
          )}
          <input
            placeholder="Password"
            id="password"
            type="text"
            name="password"
            required
            className={failed.type === "password" ? "failed" : ""}
            defaultValue={form.password}
          />
          {failed.type === "password" && (
            <span className="failed-msg">{failed.msg}</span>
          )}
          {failed.type !== "" &&
            failed.type !== "email" &&
            failed.type != "password" && (
              <div className="failed-general">{failed.msg}</div>
            )}
        </div>

        <Link className="alter-link" to="/signup">
          Sign up?
        </Link>

        <button ref={btnRef}>Login</button>
      </form>
    </div>
  );
}
