package ws

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/config"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/service/task"
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
	mu.Lock()
	for _, conn := range WsClients {
		conn.Close()
	}
	WsClients = nil
	mu.Unlock()
}

func reconnect(host string) {
	mu.Lock()
	delete(WsClients, host)
	mu.Unlock()
	err := connect(host)
	if err == nil {
		logrus.Info("connect success:", host)
		go listen(host)
		return
	}
	SendAgentChangeToExplorer()
	go func() {
		tricker := time.NewTicker(time.Duration(config.C.HeartBeatInterval) * time.Second)
		for range tricker.C {
			err = connect(host)
			if err == nil {
				tricker.Stop()
				logrus.Info("reconnect success:", host)
				return
			}
		}
	}()
}

func connect(host string) error {
	logrus.Info("connecting to ", host)
	ws := websocket.Dialer{}
	conn, _, err := ws.Dial(host, http.Header{
		"Origin":        []string{agentConfig.C.Agent},
		"Authorization": []string{"AGENT_ANALYTICS_TOKEN"},
	})
	if err != nil {
		logrus.Errorf("connect to %s error: %v", host, err)
		return err
	}
	mu.Lock()
	WsClients[host] = conn
	mu.Unlock()
	return nil
}

func listen(host string) {
	conn := WsClients[host]
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logrus.Error("read message error: ", err)
			reconnect(host)
			return
		}
		res := types.Ws_Message{}
		err = json.Unmarshal(msg, &res)
		if err != nil {
			logrus.Errorf("unmarshal message error: %v", err)
			continue
		}
		switchRoute(&res)
	}
}
func switchRoute(res *types.Ws_Message) {
	switch res.Body.MsgType {
	case types.Ws_Message_Type_Task:
		go task.HandleAnalyticsTask(res, nil)
	default:
	}
}
