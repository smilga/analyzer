package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/smilga/analyzer/api"
)

var (
	ErrNilConns = errors.New("Error conns map is nil")
)

// TODO handle closed socket connections remove them from map

type MsgType string

const (
	PingMsg MsgType = "ping"
	ConnMsg MsgType = "conn"
	CommMsg MsgType = "comm"
)

type Msg struct {
	Type    MsgType     `json:"type"`
	Message interface{} `json:"message"`
	UserID  api.UserID  `json:"userId"`
}

type Messanger struct {
	Conns map[api.UserID][]*websocket.Conn
	mutex sync.Mutex
}

func (m *Messanger) UsersOnline() []api.UserID {
	ids := []api.UserID{}
	for id, _ := range m.Conns {
		ids = append(ids, id)
	}
	return ids
}

func (m *Messanger) ReadMessage(conn *websocket.Conn) error {
	msgType, msg, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	switch msgType {
	case websocket.TextMessage:
		return m.handleTextMsg(msg, conn)
	case websocket.CloseMessage:
		return m.handleCloseMsg(string(msg), conn)
	case websocket.PingMessage:
		return m.handlePingMsg(string(msg), conn)
	}

	return nil
}

func (m *Messanger) handleTextMsg(msg []byte, conn *websocket.Conn) error {
	var message Msg

	err := json.Unmarshal(msg, &message)
	if err != nil {
		return err
	}

	switch message.Type {
	case PingMsg:
		return m.Send(conn, &Msg{PingMsg, "pong", message.UserID})
	case CommMsg:
		return m.Send(conn, &Msg{CommMsg, "any problems?", message.UserID})
	case ConnMsg:
		return m.addConn(conn, message.UserID)
	}

	return nil
}

func (m *Messanger) removeConnn(conn *websocket.Conn, id api.UserID) error {
	for i, c := range m.Conns[id] {
		if c == conn {
			m.Conns[id] = append(m.Conns[id][:i], m.Conns[id][i+1:]...)
		}
	}
	return nil
}

func (m *Messanger) addConn(conn *websocket.Conn, id api.UserID) error {
	if m.Conns == nil {
		return ErrNilConns
	}

	conn.SetCloseHandler(func(code int, text string) error {
		return m.removeConnn(conn, id)
	})

	m.Conns[id] = append(m.Conns[id], conn)

	return nil
}

func (m *Messanger) Send(conn *websocket.Conn, msg *Msg) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	ms, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := conn.WriteMessage(websocket.TextMessage, ms); err != nil {
		return err
	}

	return nil
}

func (m *Messanger) SendToUser(id api.UserID, msg *Msg) error {
	if conns, ok := m.Conns[id]; ok {
		for _, conn := range conns {
			if err := m.Send(conn, msg); err != nil {
				continue
			}
		}
	}

	return nil
}

func (m *Messanger) Broadcast(msg *Msg) error {
	for _, conns := range m.Conns {
		for _, conn := range conns {
			err := m.Send(conn, msg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *Messanger) handleCloseMsg(msg string, conn *websocket.Conn) error {
	fmt.Println("close message recieved", msg, conn)
	return nil
}

func (m *Messanger) handlePingMsg(msg string, conn *websocket.Conn) error {
	return nil
}

func NewMessanger() *Messanger {
	return &Messanger{
		make(map[api.UserID][]*websocket.Conn),
		sync.Mutex{},
	}
}
