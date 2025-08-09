package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type       string // "join", "chat", "leave"
	RoomID     string // Room name
	Sender     string // Username
	Content    string // Message content
	MaxSeat    int
	TargetUser string
}

func main() {
	// Ask for username & room
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	username = trimNewline(username)

	fmt.Print("Enter room name: ")
	roomName, _ := reader.ReadString('\n')
	roomName = trimNewline(roomName)

	serverURL := fmt.Sprintf("ws://localhost:8080/ws/%s/%s", roomName, username)

	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()

	// Handle Ctrl+C
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Listen for incoming messages
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			fmt.Printf("\n%s\nYou: ", string(msg))
		}
	}()

	// Send join event
	joinmsg := username + " Join " + roomName + " Room"
	if err := conn.WriteMessage(websocket.TextMessage, []byte(joinmsg)); err != nil {
		log.Println("Join Error: ")
	}

	// Chat input loop
	fmt.Print("You: ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()

		msg := username + "(" + roomName + "): " + text
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			log.Println("Write Error:", err)
			break
		}

		fmt.Print("You: ")
	}

	// Send leave event when exiting
	sendJSON(conn, Message{
		Type:   "leave",
		RoomID: roomName,
		Sender: username,
	})
}

func sendJSON(conn *websocket.Conn, msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("Marshal error:", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Println("Write error:", err)
	}
}

func trimNewline(s string) string {
	if len(s) > 0 && s[len(s)-1] == '\n' {
		return s[:len(s)-1]
	}
	if len(s) > 0 && s[len(s)-1] == '\r' {
		return s[:len(s)-1]
	}
	return s
}
