import { useState, useEffect, useRef, useCallback } from "react";
import { Send } from "@mui/icons-material";
import { Card, CardContent, TextField, Button, AppBar, Toolbar, Typography, Box } from "@mui/material";
import { useAuth } from "../auth/useAuth";
import { fetchApiToJson } from "../utils/apiFetch";
import { fetchPastChatMessagesByGroupId } from "../api/chat";

interface Message {
    text: string;
    sender: string;
    timestamp: string;
}

interface Module {
    id: string;
    name: string;
}

interface User {
    id: string;
    name: string;
    modules: Module[];
}

function formatDate(date: Date): string {
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const year = date.getFullYear();
    let hours = date.getHours();
    const minutes = String(date.getMinutes()).padStart(2, '0');

    const amPm = hours >= 12 ? 'PM' : 'AM';
    hours = hours % 12 || 12;

    return `${day}/${month}/${year} ${hours}:${minutes} ${amPm}`;
}

export default function ChatPage() {
    const [messages, setMessages] = useState<Message[]>([]);
    const [input, setInput] = useState("");
    const chatRef = useRef<HTMLDivElement>(null);
    const currentPath = location.toString().split("/");
    const chatID = currentPath[currentPath.length - 1];
    const ws = useRef<WebSocket | null>(null);
    const { token, user } = useAuth();
    const [userDetails, setUserDetails] = useState<User | null>(null);

    const fetchUserDetails = useCallback(async () => {
        if (!token || !user?.uid) return;
        try {
            const data = await fetchApiToJson<User>(`/user/${user.uid}`, token);
            console.log("Fetched user details:", data);
            setUserDetails(data);
        } catch (error) {
            console.error("Error fetching user details:", error);
        }
    }, [token, user]);

    const fetchPastMessages = useCallback(async () => {
        const pastChatMessages = await fetchPastChatMessagesByGroupId(token, Number(chatID));
        setMessages(() => pastChatMessages.map(message => ({
            text: message.text,
            sender: message.userName,
            timestamp: formatDate(message.timestamp)
        })))
    }, [token, chatID])

    const getTime = () => new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });

    useEffect(() => {
        void fetchUserDetails();
        void fetchPastMessages();
    }, [token, user, fetchUserDetails, fetchPastMessages]);

    useEffect(() => {
        if (ws.current?.readyState !== WebSocket.OPEN) {
            ws.current = new WebSocket(`ws://localhost:8080/api/chat/${chatID}`, ["auth", token ?? "no-token"]);

            ws.current.onopen = () => {
                console.log("WebSocket connected");
            };

            ws.current.onmessage = (event) => {
                try {
                    const msg = JSON.parse(event.data as string) as Message;
                    setMessages((prev) => [...prev, msg]);
                } catch (err) {
                    console.error("Error parsing message:", err);
                }
            };

            ws.current.onerror = (error) => {
                console.error("WebSocket error:", error);
            };

            ws.current.onclose = () => {
                console.log("WebSocket closed");
            };

            const wsCurrent = ws.current;

            return () => {
                wsCurrent.close();
            };
        }
    }, [chatID, token]);

    const currentSender = userDetails?.name ?? "You";


    const sendMessage = () => {
        if (input.trim()) {
            const newMessage: Message = {
                text: input,
                sender: currentSender,
                timestamp: getTime()
            };
            if (ws.current && ws.current.readyState === WebSocket.OPEN) {
                ws.current.send(JSON.stringify(newMessage));
            } else {
                setMessages((prevMessages) => [...prevMessages, newMessage]);
            }
            setInput("");
        }
    };

    useEffect(() => {
        if (chatRef.current) {
            chatRef.current.scrollTop = chatRef.current.scrollHeight;
        }
    }, [messages]);

    return (
        <div style={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            width: "100%",
            height: "100%",
            backgroundColor: "#ffffff",
            padding: "16px",
            boxSizing: "border-box"
        }}>
            <Card style={{
                display: "flex",
                flexDirection: "column",
                width: "100%",
                height: "100%",
                boxShadow: "0 4px 12px rgba(0, 0, 0, 0.1)",
                borderRadius: "16px",
                overflow: "hidden"
            }}>
                <AppBar position="static" color="default" style={{
                    backgroundColor: "#3b5998",
                    color: "white",
                    borderTopLeftRadius: "16px",
                    borderTopRightRadius: "16px"
                }}>
                    <Toolbar>
                        <Typography variant="h6">Chat Room</Typography>
                    </Toolbar>
                </AppBar>
                <CardContent ref={chatRef} style={{
                    flex: 1,
                    overflowY: "auto",
                    padding: "16px",
                    backgroundColor: "#ffffff"
                }}>
                    {messages.length === 0 ? (
                        <p style={{ color: "#9e9e9e", textAlign: "center" }}>Start your Chat!</p>
                    ) : (

                        messages.map((msg, index) => (
                            <div
                                key={index}
                                style={{
                                    display: "flex",
                                    flexDirection: "column",
                                    // Compare with currentSender instead of literal "You"
                                    alignItems: msg.sender === currentSender ? "flex-end" : "flex-start",
                                    marginBottom: "8px",
                                }}
                            >
                                <Typography variant="caption" style={{ color: "#555", fontWeight: "bold", marginBottom: "2px" }}>
                                    {msg.sender}
                                </Typography>
                                <Box
                                    style={{
                                        padding: "10px 14px",
                                        borderRadius: "8px",
                                        maxWidth: "60%",
                                        // Use blue background for currentSender messages.
                                        backgroundColor: msg.sender === currentSender ? "#3b5998" : "#ffffff",
                                        color: msg.sender === currentSender ? "#fff" : "#000",
                                        boxShadow: "0px 2px 4px rgba(0,0,0,0.1)",
                                    }}
                                >
                                    {msg.text}
                                </Box>
                                <Typography variant="caption" style={{ marginTop: "4px", color: "#666" }}>
                                    {msg.timestamp}
                                </Typography>
                            </div>
                        ))
                    )}
                </CardContent>
                <div style={{
                    display: "flex",
                    alignItems: "center",
                    padding: "12px",
                    borderTop: "1px solid #eee",
                    backgroundColor: "#fff",
                    borderBottomLeftRadius: "16px",
                    borderBottomRightRadius: "16px"
                }}>
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
