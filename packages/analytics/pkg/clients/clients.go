package clients

import "github.com/gorilla/websocket"

var WsClients map[string]*websocket.Conn = make(map[string]*websocket.Conn)

func GetClientByHost(host string) *websocket.Conn {
	return WsClients[host]
}
