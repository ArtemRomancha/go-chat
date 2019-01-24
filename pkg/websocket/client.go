package websocket

import (
	"github.com/ArtemRomancha/websocket-chat/pkg/chat"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type User struct {
	Username string
	Clients  map[*Client]bool
	Pool     *Pool
}

type Client struct {
	User *User
	Conn *websocket.Conn
}

type NewClient struct {
	Username string
	Conn     *websocket.Conn
}

func (c *Client) Read() {
	defer func() {
		c.User.Pool.Unregister <- c
		if err := c.Conn.Close(); err != nil {
			log.Error(err)
		}
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()

		//Connection closed from client
		if c, k := err.(*websocket.CloseError); k && c.Code == websocket.CloseNormalClosure {
			break
		}

		if err != nil {
			log.Error(err)
			return
		}
		message := chat.Message{Type: messageType, Sender: c.User.Username, Body: string(p)}
		c.User.Pool.Broadcast <- message
		log.Debugf("User %s send message: '%s'", c.User.Username, p)
	}
}
