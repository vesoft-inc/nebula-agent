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
	WsClients[host] = conn
	mu.Unlock()
	SendPendingMsg(host)
}

func SendPendingMsg(host string) {
	logrus.Infof("start send pending msg to explorer: %v", PendingMsgMap[host])
	for {
		mu.Lock()
		if len(PendingMsgMap[host]) == 0 {
			mu.Unlock()
			break
		}
		// pop first ,if send error again,reset to PendingMsgMap
		popMsg := PendingMsgMap[host][0]
		PendingMsgMap[host] = PendingMsgMap[host][1:]
		mu.Unlock()
		err := SendJsonMessage(host, popMsg)
		if err != nil {
			break
		}
	}
	delete(PendingMsgMap, host)

}

func SendJsonMessage(host string, message any) error {
	defer func() {
		if err := recover(); err != nil {
			AppendPendingMsg(host, message)
			logrus.Warnf("send msg to explorer error: %v", err)
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
		logrus.Warn(err)
	}
	return err
}

func AppendPendingMsg(host string, message any) {
	mu.Lock()
	defer mu.Unlock()
	PendingMsgMap[host] = append(PendingMsgMap[host], message)
	logrus.Warn("append pending msg to explorer: ", PendingMsgMap[host])
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
	delete(PendingMsgMap, host)
}
