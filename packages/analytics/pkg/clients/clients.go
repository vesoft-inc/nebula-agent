package clients

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var mu sync.RWMutex

var WsClients map[string]*websocket.Conn = make(map[string]*websocket.Conn)
var PendingMsgMap map[string][]any = make(map[string][]any)

func GetClientByHost(host string) *websocket.Conn {
	mu.RLock()
	defer mu.RUnlock()
	return WsClients[host]
}

func AddClientByHost(host string, conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	WsClients[host] = conn
	if PendingMsgMap[host] != nil {
		for _, msg := range PendingMsgMap[host] {
			SendJsonMessage(host, msg)
		}
		delete(PendingMsgMap, host)
	}
}

func SendJsonMessage(host string, message any) error {
	defer func() {
		if err := recover(); err != nil {
			AppendPendingMsg(host, message)
			logrus.Errorf("send msg to explorer error: %v", err)
		}
	}()
	err := func() error {
		client := GetClientByHost(host)
		if client != nil {
			return client.WriteJSON(message)
		}

		return fmt.Errorf("send msg to explorer error: %v", "conn is nil")
	}()
	if err != nil {
		AppendPendingMsg(host, message)
		logrus.Error(err)
	}
	return err
}

func AppendPendingMsg(host string, message any) {
	mu.Lock()
	defer mu.Unlock()
	PendingMsgMap[host] = append(PendingMsgMap[host], message)
}

func Clear() {
	mu.Lock()
	defer mu.Unlock()
	for _, conn := range WsClients {
		conn.Close()
	}
	WsClients = make(map[string]*websocket.Conn)
	PendingMsgMap = make(map[string][]any)
}

func DeleteClientByHost(host string) {
	mu.Lock()
	defer mu.Unlock()
	delete(WsClients, host)
}
