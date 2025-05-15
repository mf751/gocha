import { Link } from "react-router-dom";
import "./styles.css";

export default function SignUp() {
  return (
    <div className="login">
      <div className="page-title">
        <h1>Sign Up</h1>
      </div>
      <form className="login-form">
        <div className="inputs">
          <input
            id="name"
            type="text"
            placeholder="name"
            name="name"
            required
          />
          <input
            id="email"
            type="email"
            placeholder="Email"
            name="email"
            required
          />
          <input
            placeholder="Password"
            id="password"
            type="text"
            name="password"
            required
          />
        </div>

        <Link className="alter-link" to="/login">
          Login?
        </Link>

        <button>Login</button>
      </form>
    </div>
  );
}
