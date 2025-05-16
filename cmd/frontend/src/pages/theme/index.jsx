import { useState } from "react";
import { FiSun, FiMoon } from "react-icons/fi";
import "./styles.css";
import { useEffect } from "react";

export default function Theme() {
  const [theme, setTheme] = useState(() => localStorage.getItem("theme"));

  if (theme === null) {
    setTheme("light");
    localStorage.setItem("theme", "light");
  }

  if (theme === "dark") {
    document.documentElement.style.setProperty(
      "--main-color",
      "var(--d-main-color)",
    );
    document.documentElement.style.setProperty("--main-bg", "var(--d-main-bg)");
    document.documentElement.style.setProperty(
      "--login-field-bg",
      "var(--d-login-field-bg)",
    );
    document.documentElement.style.setProperty(
      "--login-field-color",
      "var(--d-login-field-color)",
    );
    document.documentElement.style.setProperty(
      "--button-color",
      "var(--d-button-color)",
    );
    document.documentElement.style.setProperty(
      "--button-bg",
      "var(--d-button-bg)",
    );
  } else {
    document.documentElement.style.setProperty(
      "--main-color",
      "var(--l-main-color)",
    );
    document.documentElement.style.setProperty("--main-bg", "var(--l-main-bg)");
    document.documentElement.style.setProperty(
      "--login-field-bg",
      "var(--l-login-field-bg)",
    );
    document.documentElement.style.setProperty(
      "--login-field-color",
      "var(--l-login-field-color)",
    );
    document.documentElement.style.setProperty(
      "--button-color",
      "var(--l-button-color)",
    );
    document.documentElement.style.setProperty(
      "--button-bg",
      "var(--l-button-bg)",
    );
  }

  function toggle() {
    setTheme((prev) => (prev === "light" ? "dark" : "light"));
    localStorage.setItem("theme", theme === "light" ? "dark" : "light");
  }

  return (
    <div className="theme">
      {theme === "light" && <FiSun className="theme-icon" onClick={toggle} />}
      {theme === "dark" && (
        <FiMoon className="theme-icon dark" onClick={toggle} />
      )}
    </div>
  );
}
