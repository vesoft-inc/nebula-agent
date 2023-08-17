package task

import (
	"testing"
	"time"

	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/types"
)

func TestGetAnalyticsTaskLog(t *testing.T) {
	InitTest()
	taskService := HandleAnalyticsTask(&types.Ws_Message{
		Body: types.Ws_Message_Body{
			Content: map[string]any{
				"action": "start",
				"task":   PageRankTask,
			},
		},
	}, host)
	for {
		if taskService.task.Status == types.TaskStatusSuccess {
			break
		}
		if taskService.task.Status == types.TaskStatusFailed {
			t.Fatalf("task failed: %s", taskService.task.JobId+"_"+taskService.task.TaskId)
			break
		}
		time.Sleep(1 * time.Second)
	}
	HandleAnalyticsTask(&types.Ws_Message{
		Body: types.Ws_Message_Body{
			Content: map[string]any{
				"action": "getLog",
				"task":   PageRankTask,
			},
		},
	}, host)
}
