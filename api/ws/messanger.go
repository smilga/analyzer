package ws

import (
	"encoding/json"
	"errors"

	"github.com/gorilla/websocket"
	"github.com/smilga/analyzer/api"
)

var (
	ErrNilConns = errors.New("Error conns map is nil")
)

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
	Conns map[api.UserID]*websocket.Conn
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

func (m *Messanger) addConn(conn *websocket.Conn, id api.UserID) error {
	if m.Conns == nil {
		return ErrNilConns
	}

	m.Conns[id] = conn
	return nil
}

func (m *Messanger) Send(conn *websocket.Conn, msg *Msg) error {
	ms, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := conn.WriteMessage(websocket.TextMessage, ms); err != nil {
		return err
	}

	return nil
}

func (m *Messanger) Broadcast(msg *Msg) error {
	for _, conn := range m.Conns {
		err := m.Send(conn, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Messanger) handleCloseMsg(msg string, conn *websocket.Conn) error {
	return nil
}

func (m *Messanger) handlePingMsg(msg string, conn *websocket.Conn) error {
	return nil
}

func NewMessanger() *Messanger {
	return &Messanger{
		make(map[api.UserID]*websocket.Conn),
	}
}
