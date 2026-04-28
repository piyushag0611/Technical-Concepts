package websocket

import (
	"net"
	"time"

	"github.com/gorilla/websocket"
)

func ConnectToServer(url string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	return conn, err
}

func SendClientMessage(conn *websocket.Conn, message string) ([]byte, error) {

	conn.WriteMessage(websocket.TextMessage, []byte(message))
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func GetServerMessages(conn *websocket.Conn) ([]string, error) {

	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	receivedMessages := []string{}
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return receivedMessages, nil
			}
		}
		receivedMessages = append(receivedMessages, string(msg))
	}
}

func CloseConnection(conn *websocket.Conn, code int) (int, error) {
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, ""))
	_, _, err := conn.ReadMessage()
	if closeErr, ok := err.(*websocket.CloseError); ok {
		return closeErr.Code, nil
	}
	return -1, err
}
