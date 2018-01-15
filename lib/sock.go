package lib

import "github.com/gorilla/websocket"

// A sock is a wrapper to a *websocket.Conn, but it also implements
// io.Writer
type sock struct {
	*websocket.Conn
}

// Write writes some bytes to the socket. The websocket type of the
// message is a TextMessage
func (s *sock) Write(data []byte) (n int, err error) {
	return len(data), s.Conn.WriteMessage(websocket.TextMessage, data)
}
