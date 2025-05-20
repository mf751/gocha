import { MdPerson } from "react-icons/md";

export default function Message({ message, isMe }) {
  if (message.message.type != 1) {
    return <div className="message transparent">{message.message.content}</div>;
  }
  if (isMe) {
    return (
      <div className="message left">
        <div className="content">
          <div className="container">
            <span className="message-text">{message.message.content}</span>
            <span className="message-time">
              {new Date(message.message.sent).toLocaleTimeString("en-GB", {
                hour: "numeric",
                minute: "2-digit",
                hour12: true,
              })}
            </span>
          </div>
        </div>
        <div className="user">
          <MdPerson className="user-icon" />
        </div>
      </div>
    );
  }
  return (
    <div className="message">
      <div className="user">
        <MdPerson className="user-icon" />
      </div>
      <div className="content">
        <span className="user-name">{message.user.name}</span>
        <div className="container">
          <span className="message-text">{message.message.content}</span>
          <span className="message-time">
            {new Date(message.message.sent).toLocaleTimeString("en-GB", {
              hour: "numeric",
              minute: "2-digit",
              hour12: true,
            })}
          </span>
        </div>
      </div>
    </div>
  );
}
