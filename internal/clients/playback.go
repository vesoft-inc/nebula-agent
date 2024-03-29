package clients

import (
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
	pb "github.com/vesoft-inc/nebula-agent/v3/pkg/proto"
)

var pbtc *PlayBackTLSConfig

type PlayBackTLSConfig struct {
	CertPath   string
	KeyPath    string
	CAPath     string
	ServerName string
	EnableSSL  bool
}

func InitPlayBackTLSConfig(caPath, certPath, keyPath, serverName string, enableSSL bool) {
	pbtc = &PlayBackTLSConfig{
		CertPath:   certPath,
		KeyPath:    keyPath,
		CAPath:     caPath,
		ServerName: serverName,
		EnableSSL:  enableSSL,
	}
}

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
	if pbtc.EnableSSL {
		cmdStr += fmt.Sprintf(" --enable_ssl=%t --cert_path=%s --key_path=%s --ca_path=%s", pbtc.EnableSSL, pbtc.CertPath, pbtc.KeyPath, pbtc.CAPath)
		if pbtc.ServerName != "" {
			cmdStr += fmt.Sprintf(" --ssl_server_SAN=%s", pbtc.ServerName)
		}
	}

	log.WithField("cmd", cmdStr).Info("Try to playback storage data...")
	cmd := exec.Command("bash", "-c", cmdStr)
	err := cmd.Run()
	if err != nil {
		log.WithError(err).Errorf("Data playback failed")
		return err
	}

	return nil
}
