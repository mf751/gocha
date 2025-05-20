import { useSelector } from "react-redux";
import "./styles.css";
import Chat from "./components/chat.jsx";

export default function Home() {
  const chats = useSelector((state) => state.chats.chats);

  if (chats === null || chats.length === 0) {
    return <h1>You Don't Have Any Chats</h1>;
  }

  return (
    <div className="chats">
      <h1>Chats</h1>
      <div className="list">
        {chats.map((obj) => (
          <Chat
            chat={obj.chat}
            key={obj.chat.id}
            lastMessage={obj.last_message}
          />
        ))}
      </div>
    </div>
  );
}
