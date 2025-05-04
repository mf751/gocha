import { useEffect, useState } from "react";

function App() {
  const [connected, setConnected] = useState(false);
  const conn = () => {
    const conn = new WebSocket("http://localhost:5050/ws");
    conn.onopen = () => setConnected(true);
  };
  return <h1 onClick={conn}>{connected ? "connected" : "disconnected"}</h1>;
}

export default App;
