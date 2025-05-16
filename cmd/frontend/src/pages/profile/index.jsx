import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";

export default function Profile() {
  const user = useSelector((state) => state.user.user);
  return <div className="profile"></div>;
}
