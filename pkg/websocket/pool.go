package websocket

import (
	"github.com/ArtemRomancha/websocket-chat/pkg/chat"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Pool struct {
	Register   chan *NewClient
	Unregister chan *Client
	Users      map[string]*User
	Broadcast  chan chat.Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *NewClient),
		Unregister: make(chan *Client),
		Users:      make(map[string]*User),
		Broadcast:  make(chan chat.Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case newClient := <-pool.Register:
			if _, ok := pool.Users[newClient.Username]; !ok {
				pool.Users[newClient.Username] = &User{
					Username: newClient.Username,
					Pool:     pool,
					Clients:  make(map[*Client]bool)}
			}
			user := pool.Users[newClient.Username]

			client := &Client{
				User: user,
				Conn: newClient.Conn,
			}
			user.Clients[client] = true
			go client.Read()

			log.Infof("User %s register new client", newClient.Username)
			if len(pool.Users[newClient.Username].Clients) == 1 {
				for _, user := range pool.Users {
					for client := range user.Clients {
						if err := client.Conn.WriteJSON(chat.Message{Type: websocket.TextMessage, Sender: "system", Body: newClient.Username + " is ONLINE"}); err != nil {
							log.Error(err)
						}
					}
				}
			}
		case client := <-pool.Unregister:
			user := pool.Users[client.User.Username]
			delete(user.Clients, client)

			log.Infof("User %s close one client", user.Username)

			if len(user.Clients) < 1 {
				userName := user.Username
				for _, user := range pool.Users {
					for client := range user.Clients {
						if err := client.Conn.WriteJSON(chat.Message{Type: websocket.TextMessage, Sender: "system", Body: userName + " is OFFLINE"}); err != nil {
							log.Error(err)
						}
					}
				}
			}
		case message := <-pool.Broadcast:
			for _, user := range pool.Users {
				for client := range user.Clients {
					if err := client.Conn.WriteJSON(message); err != nil {
						log.Error(err)
					}
				}
			}
		}
	}
}
