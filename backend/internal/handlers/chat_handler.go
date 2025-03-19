package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/services"
	"gdp8-backend/internal/utils"
)

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
	CheckOrigin:  func(_ *http.Request) bool { return true },
	Subprotocols: []string{"auth"},
}

type ChatHandler struct {
	hub         *ChatHub
	chatService services.ChatService
}

func NewChatHandler(hub *ChatHub, chatService services.ChatService) *ChatHandler {
	return &ChatHandler{
		hub:         hub,
		chatService: chatService,
	}
}

func (h *ChatHandler) ServeWs(w http.ResponseWriter, r *http.Request) {
	chatID := r.PathValue("chatID")
	studyGroupID, err := utils.ConvertToType[models.StudyGroupID](chatID)
	if chatID == "" || err != nil {
		http.Error(w, "Invalid chatID", http.StatusBadRequest)
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("DEBUG: WebSocket upgrade error: %v", err)
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}
	log.Printf("DEBUG: WebSocket upgrade successful for chatID: %s", chatID)

	h.hub.register <- RoomRegistration{room: chatID, conn: conn}
	log.Printf("DEBUG: Registered connection for chatID: %s", chatID)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("DEBUG: WebSocket read error: %v", err)
			break
		}
		log.Printf("DEBUG: ChatID %s received message: %s", chatID, message)
		h.hub.broadcast <- RoomMessage{room: chatID, message: message, messageType: messageType}

		go func() {
			var parsedMessage struct {
				Text      string `json:"text"`
				Sender    string `json:"sender"`
				Timestamp string `json:"timestamp"`
			}
			if err := json.Unmarshal(message, &parsedMessage); err != nil {
				log.Printf("Error parsing message JSON: %v", err)
			}

			err := h.chatService.InsertMessage(studyGroupID, &models.ChatMessageDetails{
				UserID:    userID,
				Text:      parsedMessage.Text,
				Timestamp: time.Now(),
			})
			if err != nil {
				log.Printf("Error inserting message: %v", err)
			}
		}()
	}

	h.hub.unregister <- RoomRegistration{room: chatID, conn: conn}
	log.Printf("DEBUG: Unregistered connection for chatID: %s", chatID)
}
