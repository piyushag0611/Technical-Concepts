package websocket_test

import (
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
	ws "websocket"

	"github.com/gorilla/websocket"
)

// Server starts and listens on a port.
func TestServerStartsAndListens(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(ws.DefaultHandler))
	defer ts.Close()

	_, err := net.Dial("tcp", ts.Listener.Addr().String())
	if err != nil {
		t.Fatalf("server not reachable: %v", err)
	}
}

// Client can connect to the server.
func TestClientCanConnect(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := ws.Upgrade(w, r)
		if err != nil {
			t.Errorf("upgrade failed: %v", err)
			return
		}
		defer conn.Close()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}))
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	clientConn, err := ws.ConnectToServer(wsURL)
	if err != nil {
		t.Fatalf("Could not connect: %v", err)
	}
	defer clientConn.Close()
}

// On connection, client receives a welcome message from the server.
func TestOnConnectionClientReceivesWelcome(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(ws.DefaultHandler))
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	clientConn, err := ws.ConnectToServer(wsURL)
	if err != nil {
		t.Fatalf("Could not connect: %v", err)
	}
	defer clientConn.Close()

	_, msg, err := clientConn.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message: %v", err)
	}
	if string(msg) != ws.WelcomeMessage {
		t.Errorf("expected %q, got %q", ws.WelcomeMessage, string(msg))
	}
}

// Client sends a message, and recevies the echoed message correctly.
func TestClientMessageReceivesEcho(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(ws.DefaultHandler))
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	clientConn, err := ws.ConnectToServer(wsURL)
	if err != nil {
		t.Fatalf("Could not connect: %v", err)
	}
	defer clientConn.Close()

	_, msg, err := clientConn.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message: %v", err)
	}
	if string(msg) != ws.WelcomeMessage {
		t.Errorf("expected %q, got %q", ws.WelcomeMessage, string(msg))
	}

	clientMessage := "PING"
	msg, err = ws.SendClientMessage(clientConn, clientMessage)
	if err != nil {
		t.Fatalf("could not read message: %v", err)
	}
	if string(msg) != clientMessage {
		t.Errorf("expected %q, got %q", ws.WelcomeMessage, string(msg))
	}
}

// Server sends an independent message, and client receives the message correctly
func TestServerMessagesClientReceives(t *testing.T) {

	serverMessage := "PING!"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := ws.Upgrade(w, r)
		if err != nil {
			t.Errorf("upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		ws.SendMessage(conn, ws.WelcomeMessage)
		ws.SendMessage(conn, serverMessage)

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}))
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	clientConn, err := ws.ConnectToServer(wsURL)
	if err != nil {
		t.Fatalf("Could not connect: %v", err)
	}
	defer clientConn.Close()

	_, msg, err := clientConn.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message: %v", err)
	}
	if string(msg) != ws.WelcomeMessage {
		t.Errorf("expected %q, got %q", ws.WelcomeMessage, string(msg))
	}

	_, msg, err = clientConn.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message: %v", err)
	}
	if string(msg) != serverMessage {
		t.Errorf("expected %q, got %q", serverMessage, string(msg))
	}
}

func TestClientCloseConnection(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(ws.DefaultHandler))
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	clientConn, err := ws.ConnectToServer(wsURL)
	if err != nil {
		t.Fatalf("Could not connect: %v", err)
	}
	defer clientConn.Close()

	_, msg, err := clientConn.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message: %v", err)
	}
	if string(msg) != ws.WelcomeMessage {
		t.Errorf("expected %q, got %q", ws.WelcomeMessage, string(msg))
	}

	closeCode, err := ws.CloseConnection(clientConn, 1000)
	if err != nil {
		t.Fatalf("error during close: %v", err)
	}
	if closeCode != 1000 {
		t.Errorf("expected close code 1000, got %d", closeCode)
	}
}

func TestServerConsecutiveMessages(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(ws.DefaultHandler))
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	clientConn, err := ws.ConnectToServer(wsURL)
	if err != nil {
		t.Fatalf("Could not connect: %v", err)
	}
	defer clientConn.Close()

	_, msg, err := clientConn.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message: %v", err)
	}
	if string(msg) != ws.WelcomeMessage {
		t.Errorf("expected %q, got %q", ws.WelcomeMessage, string(msg))
	}

	messages := []string{"hello", "how", "are", "you"}

	received := []string{}

	for _, clientMessage := range messages {
		msg, err = ws.SendClientMessage(clientConn, clientMessage)
		if err != nil {
			t.Fatalf("could not read message: %v", err)
		}
		received = append(received, string(msg))
	}

	for ind := range messages {
		if messages[ind] != received[ind] {
			t.Errorf("expected %q, got %q", messages[ind], received[ind])
		}
	}

	closeCode, err := ws.CloseConnection(clientConn, 1000)
	if err != nil {
		t.Fatalf("error during close: %v", err)
	}
	if closeCode != 1000 {
		t.Errorf("expected close code 1000, got %d", closeCode)
	}
}

func TestClientSequentialMessages(t *testing.T) {

	serverMessages := []string{"I", "am", "doing", "fine!", "thank", "you"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := ws.Upgrade(w, r)
		if err != nil {
			t.Errorf("upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		ws.SendMessage(conn, ws.WelcomeMessage)

		for _, msg := range serverMessages {
			ws.SendMessage(conn, msg)
		}

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}))
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	clientConn, err := ws.ConnectToServer(wsURL)
	if err != nil {
		t.Fatalf("Could not connect: %v", err)
	}
	defer clientConn.Close()

	_, msg, err := clientConn.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message: %v", err)
	}
	if string(msg) != ws.WelcomeMessage {
		t.Errorf("expected %q, got %q", ws.WelcomeMessage, string(msg))
	}

	recvMessages, err := ws.GetServerMessages(clientConn)
	if err != nil {
		t.Fatalf("error receiving messages: %v", err)
	}
	if !reflect.DeepEqual(serverMessages, recvMessages) {
		t.Errorf("expected %v, got %v", serverMessages, recvMessages)
	}
}

func TestMultipleClients(t *testing.T) {

	type clientMessage struct {
		conn    *websocket.Conn
		message string
	}

	ts := httptest.NewServer(http.HandlerFunc(ws.DefaultHandler))
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")

	clientConn1, err := ws.ConnectToServer(wsURL)
	if err != nil {
		t.Fatalf("Could not connect: %v", err)
	}
	defer clientConn1.Close()

	clientConn2, err := ws.ConnectToServer(wsURL)
	if err != nil {
		t.Fatalf("Could not connect: %v", err)
	}
	defer clientConn2.Close()

	clientConn3, err := ws.ConnectToServer(wsURL)
	if err != nil {
		t.Fatalf("Could not connect: %v", err)
	}
	defer clientConn2.Close()

	clients := []*websocket.Conn{clientConn1, clientConn2, clientConn3}
	for _, conn := range clients {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("could not read message: %v", err)
		}
		if string(msg) != ws.WelcomeMessage {
			t.Errorf("expected %q, got %q", ws.WelcomeMessage, string(msg))
		}
	}

	pairs := []clientMessage{
		{clientConn1, "Ping1"},
		{clientConn2, "Ping2"},
		{clientConn3, "Ping3"},
	}

	var wg sync.WaitGroup
	for _, pair := range pairs {
		wg.Add(1)
		go func(p clientMessage) {
			defer wg.Done()
			msg, err := ws.SendClientMessage(p.conn, p.message)
			if err != nil {
				t.Errorf("could not read message: %v", err)
			}
			if string(msg) != p.message {
				t.Errorf("expected %q, got %q", p.message, string(msg))
			}
		}(pair)
	}
	wg.Wait()
}
