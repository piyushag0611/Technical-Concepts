package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

const WelcomeMessage = "Welcome!"

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrade(w, r)
	if err != nil {
		return
	}
	defer conn.Close()
	conn.WriteMessage(websocket.TextMessage, []byte(WelcomeMessage))
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func SendMessage(conn *websocket.Conn, message string) {
	conn.WriteMessage(websocket.TextMessage, []byte(message))
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return upgrader.Upgrade(w, r, nil)
}
