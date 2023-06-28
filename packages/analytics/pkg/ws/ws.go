package ws

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/config"
	agentConfig "github.com/vesoft-inc/nebula-agent/v3/pkg/config"
)

var mu sync.Mutex
var WsClients map[string]*websocket.Conn

func InitWsConnect() {
	WsClients = make(map[string]*websocket.Conn)
	for _, host := range config.C.ExplorerHosts {
		connect(host)
	}
	go StartHeartBeat()
}

func CloseWsConnect() {
	StopHeartBeat()
	for _, conn := range WsClients {
		conn.Close()
	}
	WsClients = nil
}

func connect(host string) {
	ws := websocket.Dialer{}
	conn, _, err := ws.Dial(host, http.Header{
		"Origin": []string{agentConfig.C.Agent},
		"Authorization": []string{"AGENT_ANALYTICS_TOKEN"},
	})
	if err != nil {
		logrus.Errorf("connect to %s error: %v", host, err)
		return
	}
	WsClients[host] = conn
	go listen(conn)
}

func listen(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logrus.Errorf("read message error: %v", err)
			return
		}
		logrus.Infof("recive msg: %v", msg)
	}
}
