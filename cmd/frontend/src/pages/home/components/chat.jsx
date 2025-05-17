import { IoChatbubbles, IoPersonCircle } from "react-icons/io5";

export default function Chat({ chat, lastMessage }) {
  return (
    <div className="chat">
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
