import { useState, useEffect, useRef } from "react";
import { Send } from "@mui/icons-material";
import { Card, CardContent, TextField, Button, AppBar, Toolbar, Typography, Box } from "@mui/material";

interface Message {
  text: string;
  sender: string;
  timestamp: string;
}

export default function ChatUI() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const chatRef = useRef<HTMLDivElement>(null);

  const getTime = () => new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });

  const sendMessage = () => {
    if (input.trim()) {
      const newMessage: Message = { text: input, sender: "You", timestamp: getTime() };
      setMessages((prevMessages) => [...prevMessages, newMessage]);
      setInput("");
    }
  };

  useEffect(() => {
    if (chatRef.current) {
      chatRef.current.scrollTop = chatRef.current.scrollHeight;
    }
  }, [messages]);

  return (
    <div style={{ display: "flex", justifyContent: "center", alignItems: "center", height: "100vh", width: "100vw", backgroundColor: "#d0e1fd" }}>
      <Card style={{ width: "80vw", height: "90vh", display: "flex", flexDirection: "column", boxShadow: "0 4px 12px rgba(0, 0, 0, 0.1)", borderRadius: "12px", overflow: "hidden" }}>
        <AppBar position="static" color="default" style={{ backgroundColor: "#3b5998", color: "white" }}>
          <Toolbar>
            <Typography variant="h6">Chat Room</Typography>
          </Toolbar>
        </AppBar>
        <CardContent ref={chatRef} style={{ flex: 1, overflowY: "auto", padding: "16px", backgroundColor: "#b3cde8" }}>
          {messages.length === 0 ? (
            <p style={{ color: "#9e9e9e", textAlign: "center" }}>Start your Chat!</p>
          ) : (
            messages.map((msg, index) => (
              <div key={index} style={{ display: "flex", flexDirection: "column", alignItems: msg.sender === "You" ? "flex-end" : "flex-start", marginBottom: "8px" }}>
                <Typography variant="caption" style={{ color: "#555", fontWeight: "bold", marginBottom: "2px" }}>
                  {msg.sender}
                </Typography>
                <Box style={{ padding: "10px 14px", borderRadius: "8px", maxWidth: "60%", backgroundColor: msg.sender === "You" ? "#3b5998" : "#ffffff", color: msg.sender === "You" ? "#fff" : "#000", boxShadow: "0px 2px 4px rgba(0,0,0,0.1)" }}>
                  {msg.text}
                </Box>
                <Typography variant="caption" style={{ marginTop: "4px", color: "#666" }}>
                  {msg.timestamp}
                </Typography>
              </div>
            ))
          )}
        </CardContent>
        <div style={{ display: "flex", alignItems: "center", padding: "12px", borderTop: "1px solid #ddd", backgroundColor: "#fff" }}>
          <TextField
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Type a message..."
            variant="outlined"
            fullWidth
            size="small"
            onKeyDown={(e) => e.key === "Enter" && sendMessage()}
            style={{ marginRight: "8px" }}
          />
          <Button onClick={sendMessage} variant="contained" style={{ backgroundColor: "#3b5998", color: "#fff" }}>
            <Send />
          </Button>
        </div>
      </Card>
    </div>
  );
}
