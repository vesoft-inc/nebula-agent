package ws

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/config"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/task"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/types"
	agentConfig "github.com/vesoft-inc/nebula-agent/v3/pkg/config"
)

var mu sync.Mutex
var WsClients map[string]*websocket.Conn

func InitWsConnect() {
	WsClients = make(map[string]*websocket.Conn)
	for _, host := range config.C.ExplorerHosts {
		reconnect(host)
	}
	go StartHeartBeat()
	SendAgentChangeToExplorer()
}

func CloseWsConnect() {
	StopHeartBeat()
	for _, conn := range WsClients {
		conn.Close()
	}
	WsClients = nil
}

func reconnect(host string) {
	err := connect(host)
	if err == nil {
		logrus.Info("connect success:",host)
		go listen(host)
		return
	}
	mu.Lock()
	delete(WsClients, host)
	mu.Unlock()
	SendAgentChangeToExplorer()
	go func(){
		tricker := time.NewTicker(time.Duration(config.C.HeartBeatInterval)* time.Second)
		for range tricker.C {
			err = connect(host)
			if err == nil {
				tricker.Stop()
				logrus.Info("reconnect success:",host)
				return
			}
		}
	}()
}

func connect(host string) error {
	ws := websocket.Dialer{}
	conn, _, err := ws.Dial(host, http.Header{
		"Origin":        []string{agentConfig.C.Agent},
		"Authorization": []string{"AGENT_ANALYTICS_TOKEN"},
	})
	if err != nil {
		logrus.Errorf("connect to %s error: %v", host, err)
		return err
	}
	WsClients[host] = conn
	return nil
}

func listen(host string) {
	conn := WsClients[host]
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logrus.Error("read message error: ", err)
			if strings.Contains(err.Error(), "close") {
				reconnect(host)
			}
			return
		}
		res := types.Ws_Message{}
		err = json.Unmarshal(msg, &res)
		if err != nil {
			logrus.Errorf("unmarshal message error: %v", err)
			continue
		}
		switch res.Body.MsgType {
		case types.Ws_Message_Type_Task:
			task.HandleAnalyticsTask(&res)
		default:
		}
	}
}
