import "./styles.css";

export default function Login() {
  return (
    <div className="login">
      <div className="page-title">
        <h1>Login</h1>
      </div>
      <form className="login-form">
        <div className="inputs">
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

        <a className="signup" href="/signup">
          Sign up?
        </a>

        <button>Login</button>
      </form>
    </div>
  );
}
