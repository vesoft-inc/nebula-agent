package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct{}

func main() {
	//make websocket server
	ws := websocket.Upgrader{}
	httpServer := http.NewServeMux()
	httpServer.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := ws.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		// handle websocket connection
		go func() {
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					return
				}
				conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("recive msg: %v", msg)))
			}
		}()
	})
	http.ListenAndServe(":8889", httpServer)
}
