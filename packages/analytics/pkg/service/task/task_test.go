package task

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/types"
)


func TestStart(t *testing.T){
	spec := map[string]string{
			"nebula_input_metad": "",
			"nebula_input_graphd_timeout": "60000",
			"nebula_input_graphd": "192.168.8.131:9669",
			"nebula_input_user": "root",
			"nebula_input_password": "***",
			"nebula_input_space": "demo_football_2022",
			"nebula_output_types": "double",
			"input": "nebula:gflags_input",
			"nebula_output_tag": "pagerank",
			"algo_name": "pagerank",
			"encoder": "distributed",
			"nebula_input_edges": "belongto,groupedin,serve",
			"nebula_output_mode": "insert",
			"nebula_output_props": "value",
			"is_directed": "true",
			"vtype": "string",
			"need_encode": "true",
			"nebula_input_edges_props": ",,",
			"nebula_input_metad_timeout": "60000",
			"processes": "1",
			"eps": "0.0001",
			"nebula_input_storaged_timeout": "60000",
			"output": "./algo_data/0/tasks/analytics_pagerank_1",
			"damping": "0.85",
			"hosts": "192.168.8.240",
			"threads": "3",
			"iterations": "10",
	}
	task := &TaskInfo{
		JobId: "0",
		TaskId: "analytics_pagerank_1",
		Spec: spec,
	}
	wsConn := &websocket.Dialer{}
	conn,_,_ := wsConn.Dial("ws://127.0.0.1:9000/nebula_ws",nil)

	taskService := HandleAnalyticsTask(&types.Ws_Message{
		Body: types.Ws_Message_Body{
			Content: map[string]interface{}{
				"action": "start",
				"task": task,
			},
		},
	},conn)
	taskService.StartAnalyticsTask()
}