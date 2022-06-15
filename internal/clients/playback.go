package clients

import (
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
	pb "github.com/vesoft-inc/nebula-agent/pkg/proto"
)

type ServicePlayBack struct {
	dir      string
	dataPath string
	metaAddr string
}

func NewPlayBack(req *pb.DataPlayBackRequest) *ServicePlayBack {
	return &ServicePlayBack{
		dir:      req.Dir,
		dataPath: req.DataPath,
		metaAddr: req.MetaAddr,
	}
}

func (p *ServicePlayBack) PlayBack() error {
	cmdStr := fmt.Sprintf("cd %s && bin/db_playback --db_path=%s --playback_meta_server=%s", p.dir, p.dataPath, p.metaAddr)
	log.WithField("cmd", cmdStr).Debug("Try to playback storage data...")
	cmd := exec.Command("bash", "-c", cmdStr)
	err := cmd.Run()
	if err != nil {
		log.WithError(err).Errorf("Data playback failed")
		return err
	}

	return nil
}
