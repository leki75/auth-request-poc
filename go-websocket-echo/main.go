package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	fmt.Println("Listening on :8000")
	err := http.ListenAndServe(":8000", h2c.NewHandler(http.HandlerFunc(handler), &http2.Server{}))
	if err != nil {
		panic(err)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(*http.Request) bool {
		return true
	},
}

func handler(wr http.ResponseWriter, req *http.Request) {
	fmt.Printf("%s | %s %s\n", req.RemoteAddr, req.Method, req.URL)
	for header, value := range req.Header {
		fmt.Printf("%s | %s %s\n", req.RemoteAddr, header, value)
	}

	connection, err := upgrader.Upgrade(wr, req, nil)
	if err != nil {
		fmt.Printf("%s | %s\n", req.RemoteAddr, err)
		return
	}

	defer connection.Close()
	fmt.Printf("%s | upgraded to websocket\n", req.RemoteAddr)

	var message []byte

	host, err := os.Hostname()
	if err == nil {
		message = []byte(fmt.Sprintf("Request served by %s", host))
	} else {
		message = []byte(fmt.Sprintf("Server hostname unknown: %s", err.Error()))
	}

	err = connection.WriteMessage(websocket.TextMessage, message)
	if err == nil {
		var messageType int

		for {
			messageType, message, err = connection.ReadMessage()
			if err != nil {
				break
			}

			if messageType == websocket.TextMessage {
				fmt.Printf("%s | txt | %s\n", req.RemoteAddr, message)
			} else {
				fmt.Printf("%s | bin | %d byte(s)\n", req.RemoteAddr, len(message))
			}

			err = connection.WriteMessage(messageType, message)
			if err != nil {
				break
			}
		}
	}

	if err != nil {
		fmt.Printf("%s | %s\n", req.RemoteAddr, err)
	}
}
