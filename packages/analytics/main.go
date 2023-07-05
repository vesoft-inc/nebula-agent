package main

import (
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/config"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/ws"
)

func OnInit() {
	config.InitConfig("./plugins/analytics/config.yaml")
	ws.InitWsConnect()
}

func OnStop() {
	ws.CloseWsConnect()
}
