import { setLoggedIn, setUser } from "../store/slices/user";

export function SetAuthInfo(obj, dispatch) {
  dispatch(setUser(obj.user));
  dispatch(setLoggedIn(true));
  localStorage.setItem("authToken", obj["authentication_token"].token);
  localStorage.setItem("expiry", obj["authentication_token"].expiry);
}

export function Logout(dispatch) {
  dispatch(setUser({}));
  dispatch(setLoggedIn(false));
  localStorage.setItem("authToken", "");
  localStorage.setItem("expiry", "");
}
