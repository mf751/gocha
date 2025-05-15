import { useDispatch } from "react-redux";
import APIURL from "../api";
import { setLoggedIn, setUser } from "../store/slices/user";

export function SetAuthInfo(obj) {
  const dispatch = useDispatch();
  dispatch(setUser(obj.user));
  dispatch(setLoggedIn(true));
  localStorage.setItem("authToken", obj.token);
  localStorage.setItem("expiry", obj.expiry);
}

function setUserInfo() {
  const dispatch = useDispatch();
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
