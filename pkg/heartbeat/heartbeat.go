package heartbeat

import (
	"context"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vesoft-inc/nebula-agent/v3/pkg/config"
	"github.com/vesoft-inc/nebula-agent/v3/pkg/proto"
	"github.com/vesoft-inc/nebula-agent/v3/pkg/server"
)

var HeartbeatURLs []string = []string{}

func SendHeartBeat(urls []string) {
	if len(urls) == 0 {
		return
	}
	for _, url := range urls {
		go func(url string) {
			ctx := context.TODO()
			c, err := server.New(ctx, url)
			if err != nil {
				log.Fatalf("Create agent client error: %v", err)
			}
			_, err = c.SendHeartbeat(ctx, &proto.SendHeartbeatRequest{})
			if err != nil {
				logrus.Errorf("send heartbeat to %v failed: %v", url, err)
				return
			}

		}(url)
	}
	logrus.Infof("send heartbeat successfully to %v", urls)
}

func StartHeartBeat() {
	SendHeartBeat(append(HeartbeatURLs, config.C.HeartBeatHosts...))
	t := time.NewTicker(time.Duration(config.C.HeartBeatInterval) * time.Second)
	for range t.C {
		SendHeartBeat(append(HeartbeatURLs, config.C.HeartBeatHosts...))
	}
}
