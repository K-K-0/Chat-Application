package chat

import "log"

type Hub struct {
	Rooms      map[string]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

type Message struct {
	RoomID  string
	Sender  string
	Content []byte
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
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
