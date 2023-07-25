package task

import (
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/config"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/types"
	agentTask "github.com/vesoft-inc/nebula-agent/v3/pkg/task"
	"github.com/vesoft-inc/nebula-agent/v3/pkg/utils"
)

var LangTimeSpec = map[string]string{
	"job_id":                        "0",
	"task_id":                       "pagerank_1",
	"algo_name":                     "pagerank",
	"processes":                     "3",
	"nebula_output_props":           "value",
	"nebula_input_user":             "root",
	"nebula_output_tag":             "pagerank",
	"nebula_output_types":           "double",
	"nebula_input_edges_props":      ",,,,,,,,,,,,,,",
	"nebula_input_graphd":           "192.168.8.131:9669",
	"nebula_input_metad_timeout":    "60000",
	"nebula_input_storaged_timeout": "60000",
	"iterations":                    "10",
	"input":                         "nebula:gflags_input",
	"nebula_input_graphd_timeout":   "60000",
	"damping":                       "0.85",
	"threads":                       "6",
	"output":                        "/home/zhuang.miao/nebula-agent/plugins/analytics/data/langTime",
	"hosts":                         "192.168.8.240",
	"nebula_output_mode":            "insert",
	"nebula_input_edges":            "CONTAINER_OF,HAS_CREATOR,HAS_INTEREST,HAS_MEMBER,HAS_MODERATOR,HAS_TAG,HAS_TYPE,IS_LOCATED_IN,IS_PART_OF,IS_SUBCLASS_OF,KNOWS,LIKES,REPLY_OF,STUDY_AT,WORK_AT",
	"need_encode":                   "true",
	"is_directed":                   "true",
	"nebula_input_space":            "sf1_1",
	"eps":                           "0.0001",
	"nebula_input_password":         "nebula",
	"nebula_input_metad":            "",
	"encoder":                       "distributed",
	"vtype":                         "int64",
}
var PageRankSpec = map[string]string{
	"job_id":                        "0",
	"task_id":                       "pagerank_1",
	"nebula_input_metad":            "",
	"nebula_input_graphd_timeout":   "60000",
	"nebula_input_graphd":           "192.168.8.131:9669",
	"nebula_input_user":             "root",
	"nebula_input_password":         "***",
	"nebula_input_space":            "demo_football_2022",
	"nebula_output_types":           "double",
	"input":                         "nebula:gflags_input",
	"nebula_output_tag":             "pagerank",
	"algo_name":                     "pagerank",
	"encoder":                       "distributed",
	"nebula_input_edges":            "belongto,groupedin,serve",
	"nebula_output_mode":            "insert",
	"nebula_output_props":           "value",
	"is_directed":                   "true",
	"vtype":                         "string",
	"need_encode":                   "true",
	"nebula_input_edges_props":      ",,",
	"nebula_input_metad_timeout":    "60000",
	"processes":                     "1",
	"eps":                           "0.0001",
	"nebula_input_storaged_timeout": "60000",
	"damping":                       "0.85",
	"output":                        "/home/zhuang.miao/nebula-agent/plugins/analytics/data/pagerank",
	"hosts":                         "192.168.8.240",
	"threads":                       "3",
	"iterations":                    "10",
}

func InitTest() {
	agentTask.PipeShellMap = make(map[string]*agentTask.StreamShell)
	logrus.SetFormatter(&logrus.TextFormatter{})
	config.C.AnalyticsPath = "/home/zhuang.miao/nebula-analytics"
}

func TestStart(t *testing.T) {
	InitTest()
	task := TaskInfo{
		JobId:  "0",
		TaskId: "pagerank_1",
		Spec:   PageRankSpec,
	}
	wsConn := &websocket.Dialer{}
	conn, _, err := wsConn.Dial("ws://192.168.8.240:9000/nebula_ws", http.Header{
		"Origin":        []string{"192.168.8.240"},
		"Authorization": []string{"AGENT_ANALYTICS_TOKEN"},
	})
	if err != nil {
		t.Error(err)
	}
	taskService := HandleAnalyticsTask(&types.Ws_Message{
		Body: types.Ws_Message_Body{
			Content: map[string]interface{}{
				"action": "start",
				"task":   task,
			},
		},
	}, conn)
	for {
		if taskService.task.Status == types.TaskStatusSuccess {
			break
		}
		if taskService.task.Status == types.TaskStatusFailed {
			t.Fatalf("task failed: %s", taskService.task.JobId+"_"+taskService.task.TaskId)
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func TestStop(t *testing.T) {
	InitTest()
	wsConn := &websocket.Dialer{}
	conn, _, err := wsConn.Dial("ws://192.168.8.240:9000/nebula_ws", http.Header{
		"Origin":        []string{"192.168.8.240"},
		"Authorization": []string{"AGENT_ANALYTICS_TOKEN"},
	})
	if err != nil {
		t.Error(err)
	}

	task := TaskInfo{
		JobId:  "0",
		TaskId: "pagerank_1",
		Spec:   LangTimeSpec,
	}
	taskServiceStart := HandleAnalyticsTask(&types.Ws_Message{
		Body: types.Ws_Message_Body{
			Content: map[string]interface{}{
				"action": "start",
				"task":   task,
			},
		},
	}, conn)

	// async stop for wait start on next loop
	go func() {
		time.Sleep(2 * time.Second)
		HandleAnalyticsTask(&types.Ws_Message{
			Body: types.Ws_Message_Body{
				Content: map[string]interface{}{
					"action": "stop",
					"task":   task,
				},
			},
		}, conn)
	}()

	for {
		if taskServiceStart.task.Status == types.TaskStatusStopped {
			pids := utils.GetPidByName(task.JobId + "_" + task.TaskId)
			if len(pids) != 0 {
				t.Fatalf("task stop failed: %s , %s", task.JobId+"_"+task.TaskId, task.Status)
			}
			break
		}
		if taskServiceStart.task.Status == types.TaskStatusFailed || taskServiceStart.task.Status == types.TaskStatusSuccess {
			t.Fatalf("task stop failed: %s , %s", taskServiceStart.task.JobId+"_"+taskServiceStart.task.TaskId, taskServiceStart.task.Status)
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}
