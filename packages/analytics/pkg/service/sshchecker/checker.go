package sshchecker

import (
	"fmt"
	"os/user"
	"path"
	"sync"
	"time"

	"github.com/appleboy/easyssh-proxy"
	"github.com/sirupsen/logrus"

	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/clients"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/config"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/types"
)

type TaskService struct {
	host       string
	checkHosts []string
	msgId      string
	messages   map[string]string
}

func HandleCheckSSH(res *types.Ws_Message, host string) *TaskService {
	hosts := res.Body.Content["hosts"].([]string)

	t := &TaskService{
		host:       host,
		checkHosts: hosts,
		msgId:      res.Header.MsgId,
	}
	t.StartCheckSSH()
	return t
}

func (t *TaskService) StartCheckSSH() {
	u, err := user.Current()
	if err != nil {
		logrus.Error("get current user error:", err)
		return
	}
	t.messages = make(map[string]string)

	wg := &sync.WaitGroup{}
	for _, host := range t.checkHosts {
		wg.Add(1)
		t.messages[host] = ""
		keyPath := path.Join(u.HomeDir, ".ssh", "id_rsa")
		sshConfig := &easyssh.MakeConfig{
			User:    u.Username,
			Server:  host,
			KeyPath: keyPath,
		}
		go func(host string, sshConfig *easyssh.MakeConfig) {
			session, client, err := sshConfig.Connect()
			defer wg.Done()
			if err != nil {
				t.messages[host] = fmt.Sprintf("connect to %s error: %v", host, err)
				return
			}
			defer client.Close()
			defer session.Close()
			algoPath := path.Join(config.C.AnalyticsPath, "scripts/run_algo.sh")
			_, err = session.Output("cat " + algoPath)
			if err != nil {
				t.messages[host] = fmt.Sprintf("can't find algo script: %v", err)
			}
		}(host, sshConfig)
	}
	wg.Wait()

	res := types.Ws_Message{
		Header: types.Ws_Message_Header{
			SendTime: time.Now().UnixMilli(),
		},
		Body: types.Ws_Message_Body{
			MsgType: types.Ws_Message_Type_Check_SSH,
			Content: map[string]interface{}{
				"hosts":  t.checkHosts,
				"result": t.messages,
			},
		},
	}
	conn := clients.GetClientByHost(t.host)
	if conn == nil {
		logrus.Errorf("get client by host %s error", t.host)
		return
	}
	err = conn.WriteJSON(res)
	if err != nil {
		logrus.Errorf("send check ssh result to explorer error: %v", err)
	}
}
