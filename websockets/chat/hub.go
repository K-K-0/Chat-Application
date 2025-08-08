package chat

import "log"

type Hub struct {
	Rooms      map[string]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
	Commands   chan Commands
	Client     map[string]map[*Client]string
}

type Message struct {
	RoomID  string
	Sender  string
	Content []byte
}

type Commands struct {
	Action     string
	Client     *Client
	RoomID     string
	MaxSeat    int
	TargetUser string
	Content    []byte
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]map[*Client]bool),
		Client:     make(map[string]map[*Client]string),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
		Commands:   make(chan Commands, 50),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case cmd := <-h.Commands:
			switch cmd.Action {
			case "create_room":
				if h.Rooms[cmd.RoomID] == nil {
					h.Rooms[cmd.RoomID] = make(map[*Client]bool)
					log.Printf("Room %s created with max seats %d", cmd.RoomID, cmd.MaxSeat)
				}

			case "join_room":
				if _, exists := h.Rooms[cmd.RoomID]; exists {
					h.Rooms[cmd.RoomID][cmd.Client] = true
					log.Printf("Client %s joined Room %s", cmd.Client.UserID, cmd.RoomID)
				}

			case "leave_room":
				if _, exists := h.Rooms[cmd.RoomID][cmd.Client]; exists {
					delete(h.Rooms[cmd.RoomID], cmd.Client)
					log.Printf("Client %s left Room %s", cmd.Client.UserID, cmd.RoomID)
				}

			case "private_message":
				// Send only to TargetUser
				for c := range h.Rooms[cmd.RoomID] {
					if c.UserID == cmd.TargetUser {
						c.Send <- cmd.Content
					}
				}
			}

		case client := <-h.Register:

			if h.Rooms[client.RoomID] == nil {
				h.Rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.Rooms[client.RoomID][client] = true
			log.Printf("client jooined Room %s", client.RoomID)

		case client := <-h.Unregister:
			if _, ok := h.Rooms[client.RoomID][client]; ok {
				delete(h.Rooms[client.RoomID], client)
				close(client.Send)
				log.Printf("client Left Room %s", client.RoomID)
			}

		case msg := <-h.Broadcast:
			clients := h.Rooms[msg.RoomID]
			for client := range clients {
				if client.UserID != msg.Sender {
					select {
					case client.Send <- msg.Content:
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}

		}
	}
}
