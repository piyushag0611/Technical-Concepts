# Scope

- web-socket connection between client and server.
- server both echoes client messages and can also independently sends messages.
- client can send messages to the server and receive them correctly.

## Test List to pass by the web socket prototype

1. Server starts and listens on a port.
2. Client can connect to the server.
3. On connection, client receives a welcome message from the server.
4. Client sends a message, and recevies the echoed message correctly.
5. Server sends an independent message, and client receives the message correctly.
6. Client can close the connection cleanly.
7. Server handles multiple sequential messages correctly.
8. Client handles multiple sequential messages correctly.
9. Server handles multiple clients independently.