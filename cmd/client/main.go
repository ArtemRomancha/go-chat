//Package with CLI chat client
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/ArtemRomancha/websocket-chat/pkg/chat"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
)

const VERSION string ="0.01"

func printMessage(message *chat.Message) {
	if message.Sender == chat.SYSTEM_MESSAGE {
		log.Printf("\x1b[33m%s\x1b[0m", message.Body)
	} else {
		log.Printf("\x1b[32m%s\x1b[0m %s", message.Sender, message.Body)
	}
}

func connectServer(username string) *websocket.Conn {
	hostname := flag.String("hostname", "localhost", "WebSocket hostname")
	port := flag.String("port", "8080", "WebSocket port")
	path := flag.String("path", "/api/ws", "WebSocket path")
	flag.Parse()
	host := *hostname + ":" + *port
	u := url.URL{Scheme: "ws", Host: host, Path: *path}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"username": {username}})
	if err != nil {
		log.Println(err)
	}
	return c
}

func getUsername() (username string) {
	fmt.Print("Please fill Your username: ")
	_, err := fmt.Scan(&username)
	if err != nil {
		log.Println(err)
	}
	return username
}

func readUserInput(msg chan string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		fmt.Printf("\x1b[1F")
		msg <- strings.Trim(input, "\n")
	}
}

func listenServer(connection *websocket.Conn, done chan struct{}, stopListen <-chan bool) {
	defer close(done)
	for {
		message := &chat.Message{}
		err := connection.ReadJSON(message)

		if err != nil {
			if _, k := err.(*websocket.CloseError); k {
				log.Print("Server close connection")
				return
			}
			if _, k := err.(*net.OpError); k && <-stopListen {
				return
			}
			log.Println(err)
			return
		}
		printMessage(message)
	}
}

func main() {
	fmt.Println("Chat Client " + VERSION)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	username := getUsername()
	connection := connectServer(username)
	if connection == nil {
		fmt.Println("Sorry, server is unavailable!")
		return
	}

	defer func() {
		if err := connection.Close(); err != nil {
			log.Println(err)
		}
		fmt.Println("\nSee You later!")
	}()

	done := make(chan struct{})
	msg := make(chan string, 1)
	stopListen := make(chan bool, 1)

	go readUserInput(msg)

	go listenServer(connection, done, stopListen)

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			stopListen <- true
			err := connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println(err)
				return
			}
			return
		case message := <-msg:
			err := connection.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
