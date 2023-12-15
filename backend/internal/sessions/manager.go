package sessions

import (
	"fmt"

	"golang.org/x/net/websocket"
)

type SessionManager struct {
	connections map[string]*websocket.Conn
}

func New() *SessionManager {
	return &SessionManager{
		connections: make(map[string]*websocket.Conn),
	}
}

func (s *SessionManager) Add(name string, conn *websocket.Conn) {
	s.connections[name] = conn
}

func (s *SessionManager) Delete(name string) {
	delete(s.connections, name)
}

func (s *SessionManager) Send(name string, content string) error {
	conn, ok := s.connections[name]
	if !ok {
		return nil
	}
	err := websocket.Message.Send(conn, content)
	if err != nil {
		return fmt.Errorf("failed to send to %s: %w", name, err)
	}
	return nil
}
