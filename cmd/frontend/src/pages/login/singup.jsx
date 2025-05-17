import { Link, useNavigate } from "react-router-dom";
import "./styles.css";
import { useState, useRef } from "react";
import APIURL from "./../../api.js";
import { SetAuthInfo } from "../../helpers/auth.js";
import { useDispatch } from "react-redux";

export default function Signup() {
  const [form, setForm] = useState({ name: "", email: "", password: "" });
  const [failed, setFailed] = useState({});
  const btnRef = useRef();
  const navigate = useNavigate();
  const dispatch = useDispatch();

  function signup(event) {
    const name = event.get("name");
    const email = event.get("email");
    const password = event.get("password");
    setForm({ name: name, email: email, password: password });
    fetch(`${APIURL}/v1/users`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        name: name,
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
        setFailed("");
        btnRef.current.style.backgroundColor = "#41dc83";
        btnRef.current.disabled = true;
        btnRef.current.textContent = "Signed up Succuessfully!";
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
            SetAuthInfo(res, dispatch);
            navigate("/");
          });
      });
  }
  return (
    <div className="login">
      <div className="page-title">
        <h1>Sign Up</h1>
      </div>
      <form action={signup} className="login-form">
        <div className="inputs">
          <input
            id="name"
            type="text"
            placeholder="Name"
            name="name"
            required
            className={failed.type === "name" ? "failed" : ""}
            defaultValue={form.name}
          />
          {failed.type === "nae" && (
            <span className="failed-msg">{failed.msg}</span>
          )}
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

        <Link className="alter-link" to="/login">
          Login?
        </Link>

        <button ref={btnRef}>Sign Up</button>
      </form>
    </div>
  );
}
