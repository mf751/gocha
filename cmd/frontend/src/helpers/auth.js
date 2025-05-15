import { useDispatch } from "react-redux";
import APIURL from "../api";
import { setLoggedIn, setUser } from "../store/slices/user";

export function SetAuthInfo(obj, dispatch) {
  dispatch(setUser(obj.user));
  dispatch(setLoggedIn(true));
  localStorage.setItem("authToken", obj["authentication_token"].token);
  localStorage.setItem("expiry", obj["authentication_token"].expiry);
}

function setUserInfo(dispatch) {
  fetch(`${APIURL}/v1/user`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ user_id: userID }),
  })
    .then((res) => res.json())
    .then((res) => {
      dispatch(setUser(res.user));
      dispatch(setLoggedIn(true));
    });
}
