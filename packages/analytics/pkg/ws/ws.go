package ws

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/clients"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/config"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/service/sshchecker"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/service/task"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/types"
	agentConfig "github.com/vesoft-inc/nebula-agent/v3/pkg/config"
)

var mu sync.Mutex

func InitWsConnect() {
	clients.WsClients = make(map[string]*websocket.Conn)
	for _, host := range config.C.ExplorerHosts {
		reconnect(host)
	}
	go StartHeartBeat()
	SendAgentChangeToExplorer()
}

func CloseWsConnect() {
	StopHeartBeat()
	mu.Lock()
	for _, conn := range clients.WsClients {
		conn.Close()
	}
	clients.WsClients = nil
	mu.Unlock()
}

func reconnect(host string) {
	mu.Lock()
	delete(clients.WsClients, host)
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
	clients.WsClients[host] = conn
	mu.Unlock()
	return nil
}

func listen(host string) {
	conn := clients.GetClientByHost(host)
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
		switchRoute(&res, host)
	}
}

func switchRoute(res *types.Ws_Message, host string) {
	switch res.Body.MsgType {
	case types.Ws_Message_Type_Task:
		go task.HandleAnalyticsTask(res, host)
	case types.Ws_Message_Type_Check_SSH:
		go sshchecker.HandleCheckSSH(res, host)
	default:
	}
}
