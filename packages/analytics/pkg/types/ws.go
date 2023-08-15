package types

type Ws_Message struct {
	Header Ws_Message_Header `json:"header"`
	Body   Ws_Message_Body   `json:"body"`
}
type Ws_Message_Header struct {
	Version  string `json:"version,omitempty"`
	MsgId    string `json:"msgId,omitempty"`
	SendTime int64  `json:"sendTime"`
}

type Ws_Message_Type string

const (
	Ws_Message_Type_Task         Ws_Message_Type = "analytics_task"
	Ws_Message_Type_Agent        Ws_Message_Type = "analytics_agent"
	Ws_Message_Type_Check_SSH    Ws_Message_Type = "analytics_check_ssh"
	Ws_Message_Type_Machine_Info Ws_Message_Type = "analytics_machine_info"
)

type Ws_Message_Body struct {
	Product string                 `json:"product"`
	MsgType Ws_Message_Type        `json:"msgType"`
	Content map[string]interface{} `json:"content"`
}
