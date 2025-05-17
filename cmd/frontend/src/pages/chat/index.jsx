import { useNavigate, useParams } from "react-router-dom";
import "./styles.css";
import { useEffect } from "react";
import APIURL from "../../api";
import { useSelector } from "react-redux";
import { useState } from "react";
import { Link } from "react-router-dom";
import { TiArrowLeftThick } from "react-icons/ti";

export default function Chat() {
  const { chatID } = useParams();
  const [messages, setMessages] = useState([]);
  const [notInChat, setNotInChat] = useState(false);
  const chats = useSelector((state) => state.chats.chats);
  const chatsLoaded = useSelector((state) => state.chats.loaded);
  const thisChat = chats.filter((obj) => obj.chat.id === chatID)[0];
  const navigate = useNavigate();

  if (chatsLoaded && chats.length === 0) navigate("/", { replace: true });

  useEffect(() => {
    if (messages.length === 0 || !chatsLoaded) return;
    if (thisChat.last_message.id !== messages[0].id) {
      setMessages((prev) => [thisChat.last_message, ...prev]);
    }
  }, [chats, messages]);
  useEffect(() => {
    (async () => {
      try {
        const res = await fetch(
          `${APIURL}/v1/chat?id=${chatID}&start=0&size=25`,
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem("authToken")}`,
              "Content-Type": "application/json",
            },
          },
        );
        const data = await res.json();
        if (!data.error) {
          setMessages(data.data);
        }
        if (data.error === "Not a membor of chat" || res.status != 200) {
          return setNotInChat(true);
        }
        if (data.error) {
          console.log(data.error);
        }
      } catch (error) {
        if (error) console.log(error);
      }
    })();
  }, []);
  if (notInChat)
    return (
      <div className="error">
        <h3>You Are Not A Member Of This Chat</h3>
        <Link to={"/"}>Return To Chats</Link>
      </div>
    );
  if (!chatsLoaded) {
    return (
      <div
        style={{
          height: "100dvh",
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
        }}
      >
        <h1>Loading</h1>
      </div>
    );
  }
  return (
    <div className="in-chat">
      <div className="header">
        <TiArrowLeftThick
          className="icon"
          onClick={() => navigate("/", { replace: true })}
        />
        <h1>{thisChat.chat.name} </h1>
        <h2>({thisChat.members})</h2>
      </div>
      <div className="messages"></div>
    </div>
  );
}
