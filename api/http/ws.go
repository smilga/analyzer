package http

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
)

type msg string

type message struct {
	Message msg         `json:"message"`
	Content interface{} `json:"content"`
	UserID  uuid.UUID   `json:"userID"`
}

var (
	connect msg = "connect"
)

type WSMessanger struct {
	sync.Mutex
	connections map[*websocket.Conn]uuid.UUID
}

func (ws *WSMessanger) handle(conn *websocket.Conn) {
	fmt.Println("HANDEL")

	for {
		var m message

		if err := websocket.JSON.Receive(conn, &m); err != nil {
			if err == io.EOF {
				ws.removeConnection(conn)
			} else {
				log.Println("json revieve error", err)
			}
			break
		}

		switch m.Message {
		case connect:
			ws.addConnection(m.UserID, conn)
		}
	}
}

func (ws *WSMessanger) addConnection(id uuid.UUID, conn *websocket.Conn) {
	fmt.Println("add new connection")
	ws.Lock()
	defer ws.Unlock()
	ws.connections[conn] = id

}

func (ws *WSMessanger) removeConnection(conn *websocket.Conn) {
	ws.Lock()
	defer ws.Unlock()

	if _, ok := ws.connections[conn]; ok {
		delete(ws.connections, conn)
	}
}

func send(conn *websocket.Conn, data interface{}) {
	if err := websocket.JSON.Send(conn, data); err != nil {
		fmt.Println("Error sending websocket data", err)
	}
}

func broadcast(conns []*websocket.Conn, data interface{}) {
	for _, conn := range conns {
		if err := websocket.JSON.Send(conn, data); err != nil {
			fmt.Println("Error sending websocket data", err)
		}
	}
}

type WSHandler struct {
	messanger *WSMessanger
}

func (h *WSHandler) Handler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	websocket.Handler(h.messanger.handle).ServeHTTP(w, r)
}

func NewWSHandler() *WSHandler {
	return &WSHandler{
		&WSMessanger{
			sync.Mutex{},
			make(map[*websocket.Conn]uuid.UUID),
		},
	}
}
