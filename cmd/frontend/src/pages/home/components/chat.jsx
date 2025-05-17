import { IoChatbubbles, IoPersonCircle } from "react-icons/io5";
import { useNavigate } from "react-router-dom";

export default function Chat({ chat, lastMessage }) {
  const navigate = useNavigate();
  return (
    <div className="chat" onClick={() => navigate(`/chat/${chat.id}`)}>
      <div className="logo">
        {chat.is_private ? (
          <IoPersonCircle className="icon" />
        ) : (
          <IoChatbubbles className="icon" />
        )}
      </div>
      <div className="info">
        <h3>{chat.name}</h3>
        <p>{lastMessage.content}</p>
        <span>{new Date(lastMessage.sent).toLocaleTimeString()}</span>
      </div>
    </div>
  );
}
