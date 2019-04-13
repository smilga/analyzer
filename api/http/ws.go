package http

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) Upgrade(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	for {
		err = h.Messanger.ReadMessage(conn)
		if err != nil {
			conn.Close()
			break
		}
	}
}
