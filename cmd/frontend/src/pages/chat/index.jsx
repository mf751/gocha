import { useNavigate, useParams } from "react-router-dom";
import "./styles.css";
import { useEffect } from "react";
import APIURL from "../../api";
import { useSelector } from "react-redux";
import { useState } from "react";
import { Link } from "react-router-dom";
import { TiArrowLeftThick } from "react-icons/ti";
import {
  IoInformationCircleOutline,
  IoInformationCircleSharp,
  IoPersonAddOutline,
  IoPersonAddSharp,
  IoSend,
} from "react-icons/io5";
import { FiPlusCircle } from "react-icons/fi";
import Message from "./components/message";
import { useRef } from "react";

export default function Chat() {
  const { chatID } = useParams();
  const [messages, setMessages] = useState([]);
  const [notInChat, setNotInChat] = useState(false);
  const chats = useSelector((state) => state.chats.chats);
  const chatsLoaded = useSelector((state) => state.chats.loaded);
  const thisChat = chats.filter((obj) => obj.chat.id === chatID)[0];
  const navigate = useNavigate();
  const [infoShown, setInfoShowen] = useState(false);
  const [userAddShown, setUserAddShown] = useState(false);
  const user = useSelector((state) => state.user.user);
  const msgRef = useRef(null);
  const formRef = useRef(null);
  const [inputMessage, setInputMessage] = useState("");
  const [loading, setLoading] = useState(true);

  if (chatsLoaded && chats.length === 0) navigate("/", { replace: true });

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
          setLoading(false);
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

  useEffect(() => {
    if (messages.length === 0 || !chatsLoaded) return;
    if (thisChat.last_message.message.id !== messages[0].message.id) {
      setMessages((prev) => [thisChat.last_message, ...prev]);
    }
  }, [chats]);

  useEffect(() => {
    if (msgRef.current != null) {
      msgRef.current.scrollTop = msgRef.current.scrollHeight;
    }
  }, [messages, chats]);

  function sendMessage() {
    if (inputMessage === "") return;
    (async () => {
      const res = await fetch(`${APIURL}/v1/message`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem("authToken")}`,
          "Content-Type": "application/json",
        },
        method: "POST",
        body: JSON.stringify({
          chat_id: thisChat.chat.id,
          content: inputMessage,
        }),
      });
    })();
  }

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
        <div className="first-group">
          <TiArrowLeftThick
            className="icon"
            onClick={() => navigate("/", { replace: true })}
          />
          <div className="title">
            <h1>{thisChat.chat.name} </h1>
            <h2>({thisChat.members})</h2>
          </div>
        </div>
        <div className="options">
          {userAddShown ? (
            <IoPersonAddSharp
              className="add-icon"
              onClick={() => setUserAddShown((prev) => !prev)}
            />
          ) : (
            <IoPersonAddOutline
              className="add-icon"
              onClick={() => setUserAddShown((prev) => !prev)}
            />
          )}
          {infoShown ? (
            <IoInformationCircleSharp
              className="info-icon"
              onClick={() => setInfoShowen((prev) => !prev)}
            />
          ) : (
            <IoInformationCircleOutline
              className="info-icon"
              onClick={() => setInfoShowen((prev) => !prev)}
            />
          )}
        </div>
      </div>
      {loading && (
        <div className="loading">
          <h1>Loading</h1>
        </div>
      )}
      <div ref={msgRef} className="messages">
        {[...messages].reverse().map((msg) => (
          <Message
            message={msg}
            isMe={msg.user.id === user.id}
            key={msg.message.id}
          />
        ))}
      </div>
      <form className="typing" action={sendMessage} ref={formRef}>
        <FiPlusCircle className="add-icon" />
        <textarea
          name="message"
          className="typing-area"
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey && e.target.value.trim()) {
              e.preventDefault();
              formRef.current.requestSubmit();
            }
          }}
          onChange={(e) => setInputMessage(e.target.value)}
        ></textarea>
        <button type="submit">
          <IoSend className="send" />
        </button>
      </form>
    </div>
  );
}
