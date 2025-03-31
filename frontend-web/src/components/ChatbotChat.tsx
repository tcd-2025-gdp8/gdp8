import { useState, useEffect, useRef } from "react";
import { Send } from "@mui/icons-material";
import { AppBar, Toolbar, Typography, Box, TextField, Button, CardContent } from "@mui/material";
import { prompt } from "../api/chatbot";
import { useAuth } from "../auth/useAuth";

interface Message {
    text: string;
    sender: string;
    timestamp: string;
}

interface ChatbotChatProps {
    onClose: () => void;
}

function extractGroupNumber(url: string): string {
    const match = /\/study-groups\/(\d+)\/files/.exec(url);
    return match ? match[1] : "";
}

export default function ChatbotChat({ onClose }: ChatbotChatProps) {
    const [botMessages, setBotMessages] = useState<Message[]>([]);
    const [botInput, setBotInput] = useState("");
    const chatID = extractGroupNumber(window.location.toString());
    const { token } = useAuth();
    const botChatRef = useRef<HTMLDivElement>(null);

    const sendBotMessage = async () => {
        if (botInput.trim()) {
            const userMessage: Message = {
                text: botInput,
                sender: "You",
                timestamp: new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }),
            };
            setBotMessages((prev) => [...prev, userMessage]);

            try {
                const response = await prompt({
                    token: token,
                    chatID: chatID,
                    message: botInput,
                    memory: "file"
                });
                const botReply: Message = {
                    text: response.message,
                    sender: "Chatbot",
                    timestamp: new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }),
                };
                setBotMessages((prev) => [...prev, botReply]);
            } catch {
                const errorReply: Message = {
                    text: "Error fetching response.",
                    sender: "Chatbot",
                    timestamp: new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }),
                };
                setBotMessages((prev) => [...prev, errorReply]);
            }

            setBotInput("");
        }
    };


    // Scroll to the bottom whenever botMessages change
    useEffect(() => {
        if (botChatRef.current) {
            botChatRef.current.scrollTop = botChatRef.current.scrollHeight;
        }
    }, [botMessages]);

    return (
        <div style={{ 
            position: "fixed", 
            bottom: "80px", 
            right: "20px", 
            width: "300px", 
            backgroundColor: "#fff", 
            borderRadius: "12px",
            boxShadow: "0px 4px 12px rgba(0,0,0,0.2)", 
            display: "flex", 
            flexDirection: "column"
        }}>
            <AppBar position="static" color="primary" style={{ borderRadius: "12px 12px 0 0", backgroundColor: "#3b5998" }}>
                <Toolbar style={{ display: "flex", justifyContent: "space-between" }}>
                    <Typography variant="h6" style={{ fontSize: "16px" }}>Chatbot</Typography>
                    <Button onClick={onClose} style={{ color: "#fff" }}>X</Button>
                </Toolbar>
            </AppBar>
            <CardContent ref={botChatRef} style={{ 
                flex: 1, 
                overflowY: "auto", 
                padding: "10px", 
                maxHeight: "300px" 
            }}>
                {botMessages.length === 0 ? (
                    <p style={{ color: "#9e9e9e", textAlign: "center" }}>Chat with the bot!</p>
                ) : (
                    botMessages.map((msg, index) => (
                        <div key={index} style={{ 
                            display: "flex", 
                            flexDirection: "column", 
                            alignItems: msg.sender === "You" ? "flex-end" : "flex-start", 
                            marginBottom: "8px" 
                        }}>
                            <Typography variant="caption" style={{ fontWeight: "bold", color: "#555" }}>
                                {msg.sender}
                            </Typography>
                            <Box style={{ 
                                padding: "10px", 
                                borderRadius: "8px", 
                                maxWidth: "80%", 
                                backgroundColor: msg.sender === "You" ? "#3b5998" : "#f0f0f0", 
                                color: msg.sender === "You" ? "#fff" : "#000"
                            }}>
                                {msg.text}
                            </Box>
                            <Typography variant="caption" style={{ color: "#666" }}>
                                {msg.timestamp}
                            </Typography>
                        </div>
                    ))
                )}
            </CardContent>
            <div style={{ display: "flex", padding: "8px", borderTop: "1px solid #eee" }}>
                <TextField 
                    value={botInput} 
                    onChange={(e) => setBotInput(e.target.value)}
                    placeholder="Ask something..."
                    variant="outlined"
                    size="small"
                    fullWidth
                    onKeyDown={(e) => { if (e.key === "Enter") void sendBotMessage(); }} 
                />
                <Button onClick={() => void sendBotMessage()} style={{ marginLeft: "8px", backgroundColor: "#3b5998", color: "#fff" }}>
                    <Send />
                </Button>
            </div>
        </div>
    );
}
