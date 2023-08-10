package task

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/clients"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/config"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/types"
	agentTask "github.com/vesoft-inc/nebula-agent/v3/pkg/task"
	"github.com/vesoft-inc/nebula-agent/v3/pkg/utils"
)

func (t *TaskService) GetAnalyticsTaskLog() error {
	id := t.task.JobId + "_" + t.task.TaskId
	analyticPath := config.C.AnalyticsPath
	if t.task.AnalyticsPath != "" {
		analyticPath = t.task.AnalyticsPath
	}
	logPath := path.Join(analyticPath, "logs", id, "all.log")
	// if agentTask.IsShellRunning(id) {
	// 	go t.StartPipeLog(id, logPath)
	// } else {
	lines, err := GetSomeLinesLogWithPath(logPath)
	if err != nil {
		logrus.Error("get 200 lines log error:", err)
		return nil
	}
	t.SendLogToExplorer(strings.Join(lines, "\n"))
	// }
	return nil
}

func (t *TaskService) StartPipeLog(id string, path string) error {
	// Start tailing the log file
	cmd := exec.Command("tail", "-f", path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Error("get stdout pipe error:", err)
		return err
	}
	if err := cmd.Start(); err != nil {
		logrus.Error("start tail error:", err)
		return err
	}

	reader := bufio.NewReader(stdout)

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			logrus.Error("read log error:", err)
			return err
		}
		if line != "" {
			t.SendLogToExplorer(line)
		}
		if err == io.EOF {
			// eof sometimes means the file is not ready yet, wait for a while
			if agentTask.IsShellRunning(id) {
				time.Sleep(100 * time.Millisecond)
				t.GetAnalyticsTaskLog()
			}
			break
		}
	}
	return nil
}

func (t *TaskService) StopPipeLog() {
	id := t.task.JobId + "_" + t.task.TaskId
	logPath := path.Join(config.C.AnalyticsPath, "logs", id, "all.log")
	pids := utils.GetPidByName(logPath)
	utils.KillProcessByPids(pids)
}

func (t *TaskService) SendLogToExplorer(text string) {
	if testing.Verbose() {
		log.Print("------->", text)
	}
	conn := clients.GetClientByHost(t.host)
	conn.WriteJSON(types.Ws_Message{
		Header: types.Ws_Message_Header{
			MsgId:    t.msgId,
			SendTime: time.Now().UnixMilli(),
		},
		Body: types.Ws_Message_Body{
			MsgType: types.Ws_Message_Type_Task,
			Content: map[string]interface{}{
				"action": "log",
				"task": map[string]interface{}{
					"jobId":  t.task.JobId,
					"taskId": t.task.TaskId,
					"data":   text,
				},
			},
		},
	})
}

func GetSomeLinesLogWithPath(path string) ([]string, error) {
	//read log
	file, err := os.Open(path)
	if err != nil {
		logrus.Error("open file error: ", err)
		return nil, err
	}
	defer file.Close()
	if err != nil {
		logrus.Error("read file error: ", err)
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	var lines []string
	i := int32(0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if i > config.C.LogNum {
			lines = append(lines[:100], lines[101:]...)
		}
		i++
	}
	return lines, nil
}
