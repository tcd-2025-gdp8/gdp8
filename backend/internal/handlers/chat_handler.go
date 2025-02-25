package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"firebase.google.com/go/v4/auth"
        "gdp8-backend/internal/middleware"
)

// RoomRegistration represents a connection joining a chat room.
type RoomRegistration struct {
	room string
	conn *websocket.Conn
}

type RoomMessage struct {
	room        string
	message     []byte
	messageType int
}

type ChatHub struct {
	rooms      map[string]map[*websocket.Conn]bool
	register   chan RoomRegistration
	unregister chan RoomRegistration
	broadcast  chan RoomMessage
}

func NewChatHub() *ChatHub {
	return &ChatHub{
		rooms:      make(map[string]map[*websocket.Conn]bool),
		register:   make(chan RoomRegistration),
		unregister: make(chan RoomRegistration),
		broadcast:  make(chan RoomMessage),
	}
}

func (h *ChatHub) Run() {
	for {
		select {
		case reg := <-h.register:
			if _, ok := h.rooms[reg.room]; !ok {
				h.rooms[reg.room] = make(map[*websocket.Conn]bool)
			}
			h.rooms[reg.room][reg.conn] = true

		case reg := <-h.unregister:
			if conns, ok := h.rooms[reg.room]; ok {
				if _, exists := conns[reg.conn]; exists {
					delete(conns, reg.conn)
					reg.conn.Close()
					if len(conns) == 0 {
						delete(h.rooms, reg.room)
					}
				}
			}

		case msg := <-h.broadcast:
			if conns, ok := h.rooms[msg.room]; ok {
				for conn := range conns {
					if err := conn.WriteMessage(msg.messageType, msg.message); err != nil {
						log.Println("Broadcast write error:", err)
						delete(conns, conn)
						conn.Close()
					}
				}
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	Subprotocols: []string{"auth", "accessToken"},
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ChatHandler struct {
	hub          *ChatHub
	firebaseAuth *auth.Client
}

func NewChatHandler(hub *ChatHub, firebaseAuth *auth.Client) *ChatHandler {
	return &ChatHandler{hub: hub, firebaseAuth: firebaseAuth}
}

func (h *ChatHandler) ServeWs(w http.ResponseWriter, r *http.Request) {
	log.Printf("DEBUG: Received headers: %+v", r.Header)

	protocolHeader := r.Header.Get("Sec-WebSocket-Protocol")
	log.Printf("DEBUG: Sec-WebSocket-Protocol header: %s", protocolHeader)
	protocols := strings.Split(protocolHeader, ",")
	for i := range protocols {
		protocols[i] = strings.TrimSpace(protocols[i])
	}
	if len(protocols) < 2 {
		http.Error(w, "Missing token in subprotocol", http.StatusUnauthorized)
		return
	}
	tokenString := protocols[1]
	log.Printf("DEBUG: Token from subprotocol: %s", tokenString)

	decodedToken, err := h.firebaseAuth.VerifyIDToken(r.Context(), tokenString)
	if err != nil {
		log.Printf("DEBUG: Token verification failed: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("DEBUG: Token verified. UID: %s", decodedToken.UID)

	ctx := context.WithValue(r.Context(), middleware.UIDCtxKey{}, decodedToken.UID)
	r = r.WithContext(ctx)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Chat ID missing in URL", http.StatusBadRequest)
		return
	}
	chatID := parts[3]
	log.Printf("DEBUG: Received chatID: %s", chatID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("DEBUG: WebSocket upgrade error: %v", err)
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}
	log.Printf("DEBUG: WebSocket upgrade successful for chatID: %s, UID: %s", chatID, decodedToken.UID)

	h.hub.register <- RoomRegistration{room: chatID, conn: conn}
	log.Printf("DEBUG: Registered connection for chatID: %s", chatID)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("DEBUG: WebSocket read error: %v", err)
			break
		}
		log.Printf("DEBUG: ChatID %s received message from UID %s: %s", chatID, decodedToken.UID, message)
		h.hub.broadcast <- RoomMessage{room: chatID, message: message, messageType: messageType}
	}

	h.hub.unregister <- RoomRegistration{room: chatID, conn: conn}
	log.Printf("DEBUG: Unregistered connection for chatID: %s", chatID)
}
