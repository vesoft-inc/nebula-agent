package task

import (
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/types"
)

type TaskInfo struct {
	JobId  string
	TaskId string
	Status TaskStatus
}
type TaskStatus string

const (
	TaskStatusRunning     TaskStatus = "running"
	TaskStatusSuccess     TaskStatus = "success"
	TaskStatusFailed      TaskStatus = "failed"
	TaskStatusInterrupted TaskStatus = "interrupted"
	TaskStatusStopping    TaskStatus = "stopping"
)

var taskMap = make(map[string]TaskInfo)

func HandleAnalyticsTask(res *types.Ws_Message) {
	action := res.Body.Content["action"].(string)
	switch action {
	case "start":
		StartAnalyticsTask(res)
	case "stop":
		StopAnalyticsTask(res)
	}
}

func StartAnalyticsTask(res *types.Ws_Message) {
	// TODO
}
func StopAnalyticsTask(res *types.Ws_Message) {
	// TODO
}
