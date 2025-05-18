import { MdPerson } from "react-icons/md";

export default function Message({ message, isMe }) {
  return (
    <div className={`message ${isMe ? "" : "left"}`}>
      <div className="user">
        <MdPerson className="user-icon" />
      </div>
      <div className="content">{message.content}</div>
    </div>
  );
}
