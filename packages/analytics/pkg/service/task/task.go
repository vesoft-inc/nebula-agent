package task

import (
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/types"
	agentTask "github.com/vesoft-inc/nebula-agent/v3/pkg/task"
)

type TaskInfo struct {
	JobId  string `json:"jobId"`
	TaskId string `json:"taskId"`
	ExecFile string `json:"exec_file"`
	Spec     map[string]string `json:"spec"`

	Status types.TaskStatus `json:"status"`
	StartTime int64 `json:"start_time"`
	EndTime int64 `json:"end_time"`
}

type TaskService struct {
	conn *websocket.Conn
	task *TaskInfo
}

func HandleAnalyticsTask(res *types.Ws_Message,conn *websocket.Conn) {
	action := res.Body.Content["action"].(string)
	task := res.Body.Content["task"].(TaskInfo)
	t := &TaskService{
		conn: conn,
		task: &task,
	}
	switch action {
	case "start":
		t.StartAnalyticsTask()
	case "stop":
		t.StopAnalyticsTask()
	}
}

func (t *TaskService) StartAnalyticsTask() {
	 taskInfo := t.task
	 cmd := task2cmd(taskInfo, true)
	 cmdWithoutPwd := task2cmd(taskInfo, false)
	 id := taskInfo.JobId + "-" + taskInfo.TaskId
	 taskInfo.Status = types.TaskStatusRunning
	 taskInfo.StartTime = time.Now().Unix()
	 t.SendTaskStatusToExplorer()
	 go func() {
		  err := agentTask.RunStreamShell(id, cmd,func(msg string) error{
				return nil
			})
			if err != nil {
				logrus.Errorf("run task %s failed: %s", cmdWithoutPwd, err)
				taskInfo.Status = types.TaskStatusFailed
			} else {
				logrus.Infof("run task %s success", cmdWithoutPwd)
				taskInfo.Status = types.TaskStatusSuccess
			}
			taskInfo.EndTime = time.Now().Unix()
			t.SendTaskStatusToExplorer()
	 }()
	 
}
func (t *TaskService) StopAnalyticsTask() {
	id := t.task.JobId + "-" + t.task.TaskId
	err := agentTask.StopStreamShell(id)
	if err != nil {
		logrus.Errorf("stop task %s failed: %s", id, err)
	}
	t.task.Status = types.TaskStatusStopped
	t.task.EndTime = time.Now().Unix()
	t.SendTaskStatusToExplorer()
}

func (t *TaskService) SendTaskStatusToExplorer() {
		t.conn.WriteJSON(types.Ws_Message{
			Header: types.Ws_Message_Header{
				SendTime: time.Now().Unix(),
			},
			Body: types.Ws_Message_Body{
				MsgType: "analytics_task",
				Content: map[string]interface{}{
					"task": map[string]interface{}{
						"jobId":  t.task.JobId,
						"taskId": t.task.TaskId,
						"status": t.task.Status,
					},
					"spendTime": t.task.EndTime - t.task.StartTime,
				},
			},
	})
}

func task2cmd(task *TaskInfo, filterPwd bool) string {
	cmd := task.ExecFile
	if jobId, exist := task.Spec["job_id"]; exist {
		cmd = cmd + " --job_id '" + jobId + "' "
	}
	if taskId, exist := task.Spec["task_id"]; exist {
		cmd = cmd + " --task_id '" + taskId + "' "
	}
	for key, value := range task.Spec {
		if key == "job_id" || key == "task_id" {
			continue
		}

		if filterPwd && strings.Contains(key, "password") {
			cmd = cmd + " --" + key + " '***' "
		} else {
			cmd = cmd + " --" + key + " '" + value + "' "
		}
	}
	return cmd
}
