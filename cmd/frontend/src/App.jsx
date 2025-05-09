import { useEffect, useState } from "react";

function App() {
  const [connected, setConnected] = useState(false);
  const [conn1, setConn1] = useState();
  const [conn2, setConn2] = useState();
  const authToken = "52V6HUN76UILQ3SUVJLMULWJO4";
  const authToken2 = "SXXTGUZMNFRC45KFSXKZTGVC24";
  const chatID = "7f9401bf-9ce4-4cd1-82aa-b65e0a722d4d";
  useEffect(() => {
    const conn = new WebSocket(`ws://localhost:5050/v1/ws?token=${authToken}`);
    conn.onopen = () => setConnected(true);
    conn.onclose = () => setConnected(false);
    conn.onmessage = (event) => console.log("From conn1: ", event);
    setConn1(conn);
    const conn2 = new WebSocket(
      `ws://localhost:5050/v1/ws?token=${authToken2}`,
    );
    conn2.onmessage = (evt) => console.log("From conn 2: ", evt);
    conn2.onopen = () => console.log("conn 2 opened");
    conn2.onclose = () => console.log("conn 2 closed");
    setConn2(conn2);
  }, []);
  function sendEvent(nun) {
    if (nun === 1) {
      conn1.send(
        JSON.stringify({
          type: "send_message",
          payload: { message: "test", chat_id: chatID },
        }),
      );
    } else {
      conn2.send(
        JSON.stringify({
          type: "send_message",
          payload: { message: "test", chat_id: chatID },
        }),
      );
    }
  }
  return (
    <div>
      <h1>{connected ? "connected" : "disconnected"}</h1>;
      <button onClick={() => sendEvent(1)}>Send from conn 1</button>
      <button onClick={() => sendEvent(2)}>Send from conn 2</button>
    </div>
  );
}

export default App;
