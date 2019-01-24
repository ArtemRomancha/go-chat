//Main package for chat server
package main

import (
	"flag"
	"github.com/ArtemRomancha/websocket-chat/pkg/logFormatter"
	"github.com/ArtemRomancha/websocket-chat/pkg/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const VERSION string ="0.01"

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	log.Debugf("Initial request to a websocket from user: %s", r.Header.Get("username"))
	conn, err := websocket.Upgrade(w, r)

	if err != nil {
		log.Error(err)
		return
	}

	pool.Register <- &websocket.NewClient{Username: r.Header.Get("username"), Conn: conn}
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/api/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

func init() {
	log.SetFormatter(new(logFormatter.CustomLogFormatter))
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Info("Chat App " + VERSION)
	setupRoutes()
	port := flag.String("port", "8080", "WebSocket port")
	flag.Parse()
	log.Error(http.ListenAndServe(":" + *port, nil))
}
