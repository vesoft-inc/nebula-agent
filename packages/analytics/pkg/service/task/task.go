package task

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/config"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/types"
	agentTask "github.com/vesoft-inc/nebula-agent/v3/pkg/task"
	"github.com/vesoft-inc/nebula-agent/v3/pkg/utils"
)

type TaskInfo struct {
	JobId         string         `json:"jobId"`
	TaskId        string         `json:"taskId"`
	Spec          map[string]any `json:"spec"`
	AnalyticsPath string         `json:"analytics_path"`

	Status    types.TaskStatus `json:"status"`
	StartTime int64            `json:"start_time"`
	EndTime   int64            `json:"end_time"`
}

type TaskService struct {
	conn  *websocket.Conn
	task  *TaskInfo
	msgId string
}

func HandleAnalyticsTask(res *types.Ws_Message, conn *websocket.Conn) *TaskService {
	action := res.Body.Content["action"].(string)
	taskContent := res.Body.Content["task"].(map[string]any)
	task := TaskInfo{
		JobId:  taskContent["jobId"].(string),
		TaskId: taskContent["taskId"].(string),
		Spec:   taskContent["spec"].(map[string]any),
		Status: types.TaskStatusRunning,
	}
	t := &TaskService{
		conn:  conn,
		task:  &task,
		msgId: res.Header.MsgId,
	}
	switch action {
	case "start":
		t.StartAnalyticsTask()
	case "stop":
		t.StopAnalyticsTask()
	case "getLog":
		t.GetAnalyticsTaskLog()
	case "stopLog":
		t.StopPipeLog()
	}
	return t
}

func (t *TaskService) StartAnalyticsTask() {
	taskInfo := t.task
	cmd := task2cmd(taskInfo, true)
	cmdWithoutPwd := task2cmd(taskInfo, false)
	id := taskInfo.JobId + "_" + taskInfo.TaskId
	taskInfo.Status = types.TaskStatusRunning
	taskInfo.StartTime = time.Now().Unix()
	t.SendTaskStatusToExplorer()
	go func() {
		logrus.Info("start task: ", cmdWithoutPwd)
		err := agentTask.RunStreamShell(id, cmd, func(msg string) error {
			return nil
		})
		if err != nil {
			if err.Error() == "stop stream shell" {
				logrus.Infof("stop task %s succeed:", id)
				taskInfo.Status = types.TaskStatusStopped
			} else {
				logrus.Errorf("run task failed:%s err: %s", cmdWithoutPwd, err)
				taskInfo.Status = types.TaskStatusFailed
			}
		} else {
			taskInfo.Status = types.TaskStatusSuccess
		}
		t.SendTaskStatusToExplorer()
	}()
}

func (t *TaskService) StopAnalyticsTask() {
	id := t.task.JobId + "_" + t.task.TaskId
	err := agentTask.StopStreamShell(id)
	if err != nil {
		logrus.Errorf("stop task %s failed: %s", id, err)
	}
	t.KillAnalyticsProcess()
}

func (t *TaskService) KillAnalyticsProcess() {
	logDirId := t.task.JobId + "_" + t.task.TaskId
	pids := utils.GetPidByName(logDirId)
	logrus.Infof("kill task %s pids: %v", logDirId, pids)
	utils.KillProcessByPids(pids)
}

func (t *TaskService) SendTaskStatusToExplorer() {
	content := map[string]any{
		"task": map[string]any{
			"jobId":     t.task.JobId,
			"taskId":    t.task.TaskId,
			"status":    t.task.Status,
			"startTime": t.task.StartTime,
			"endTime":   time.Now().Unix(),
		},
	}
	logrus.Info("send task status to explorer: ", t.task.JobId, "_", t.task.TaskId, " status: ", t.task.Status)
	t.conn.WriteJSON(types.Ws_Message{
		Header: types.Ws_Message_Header{
			SendTime: time.Now().Unix(),
			MsgId:    t.msgId,
		},
		Body: types.Ws_Message_Body{
			MsgType: types.Ws_Message_Type_Task,
			Content: content,
		},
	})
}

func task2cmd(task *TaskInfo, filterPwd bool) string {
	cmd := path.Join(config.C.AnalyticsPath, "scripts/run_algo.sh")
	if task.AnalyticsPath != "" {
		cmd = path.Join(task.AnalyticsPath, "scripts/run_algo.sh")
	}
	cmd = cmd + " --job_id '" + task.JobId + "' "
	cmd = cmd + " --task_id '" + task.TaskId + "' "
	for key, value := range task.Spec {
		if key == "job_id" || key == "task_id" {
			continue
		}

		if filterPwd && strings.Contains(key, "password") {
			cmd = cmd + " --" + key + " '***' "
		} else {
			cmd = cmd + " --" + fmt.Sprintf("%v", key) + " '" + fmt.Sprintf("%v", value) + "' "
		}
	}
	return cmd
}
