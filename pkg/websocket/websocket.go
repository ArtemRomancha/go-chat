package websocket

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func Upgrade(responseWriter http.ResponseWriter, request *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(responseWriter, request, nil)
	if err != nil {
		log.Error(err)
	}
	return ws, err
}
