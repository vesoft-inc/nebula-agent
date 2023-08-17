package sshchecker

import (
	"testing"

	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/types"
)

func TestHandleCheckSSH(t *testing.T) {
	// create a test message
	testMsg := types.Ws_Message{
		Header: types.Ws_Message_Header{
			MsgId: "test-msg-id",
		},
		Body: types.Ws_Message_Body{
			Content: map[string]interface{}{
				"hosts": []string{"192.168.8.240", "host2", "host3"},
			},
		},
	}

	// create a test TaskService
	testHost := "test-host"
	testTask := HandleCheckSSH(&testMsg, testHost)

	// check that the TaskService was created correctly
	if testTask.host != testHost {
		t.Errorf("Expected host %s, but got %s", testHost, testTask.host)
	}
	if len(testTask.checkHosts) != 3 {
		t.Errorf("Expected 3 check hosts, but got %d", len(testTask.checkHosts))
	}
	if testTask.msgId != "test-msg-id" {
		t.Errorf("Expected msgId test-msg-id, but got %s", testTask.msgId)
	}
	if len(testTask.messages) != len(testMsg.Body.Content["hosts"].([]string)) {
		t.Errorf("Expected messages map, but got nil")
	}
	t.Log("TestHandleCheckSSH passed", testTask.messages)
}
