import { useEffect, useState } from "react";

function App() {
  const [connected, setConnected] = useState(false);
  const [conn, setConn] = useState(null);
  const authToken = "52V6HUN76UILQ3SUVJLMULWJO4";
  const authToken2 = "SXXTGUZMNFRC45KFSXKZTGVC24";
  useEffect(() => {
    const conn = new WebSocket(`ws://localhost:5050/v1/ws?token=${authToken}`);
    conn.onopen = () => setConnected(true);
    conn.onclose = () => setConnected(false);
    conn.onmessage = (event) => console.log(event);
    setConn(conn);
  }, []);
  function sendEvent() {
    conn.send(
      JSON.stringify({
        type: "send_message",
        payload: { message: "test", chat_id: chatID },
      }),
    );
  }
  return (
    <div>
      <h1>{connected ? "connected" : "disconnected"}</h1>;
      <button onClick={sendEvent}>Send</button>
    </div>
  );
}

export default App;
