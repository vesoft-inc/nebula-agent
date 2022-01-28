package clients

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	pb "github.com/vesoft-inc/nebula-agent/pkg/proto"
)

type ServiceName string

const (
	ServiceName_All      ServiceName = "all"
	ServiceName_Metad    ServiceName = "metad"
	ServiceName_Storaged ServiceName = "storaged"
	ServiceName_Graphd   ServiceName = "graphd"
	ServiceName_Unknown  ServiceName = "unknown"
)

func toName(r pb.ServiceRole) ServiceName {
	switch r {
	case pb.ServiceRole_ALL:
		return ServiceName_All
	case pb.ServiceRole_META:
		return ServiceName_Metad
	case pb.ServiceRole_STORAGE:
		return ServiceName_Storaged
	case pb.ServiceRole_GRAPH:
		return ServiceName_Graphd
	default:
		log.Errorf("Bad format role string: %s", r.String())
	}
	return ServiceName_Unknown
}

// Service keeps the metad/storaged/graphd's role name and root dir which
// should have scripts to start/stop service
type Service struct {
	name ServiceName
	dir  string
}

func FromStartReq(req *pb.StartServiceRequest) *Service {
	return &Service{
		name: toName(req.GetRole()),
		dir:  req.GetDir(),
	}
}

func FromStopReq(req *pb.StopServiceRequest) *Service {
	return &Service{
		name: toName(req.GetRole()),
		dir:  req.GetDir(),
	}
}

func FromStatusReq(req *pb.ServiceStatusRequest) *Service {
	return &Service{
		name: toName(req.GetRole()),
		dir:  req.GetDir(),
	}
}

// ServiceDaemon will start/stop metad/storaged/graphd in the service machine
// through scripts providing by the nebula
type ServiceDaemon struct {
	s *Service
}

func NewDaemon(s *Service) (*ServiceDaemon, error) {
	if s == nil || s.dir == "" {
		return nil, fmt.Errorf("%s's is not found or root dir is empty", s.name)
	}

	d := &ServiceDaemon{
		s: s,
	}
	return d, nil
}

// TODO(spw): The code will only check the scripts return code. And when return code is 0,
// it does not mean the service has must been started successfully. Then we need to
// check the service process status by other means in the future
func (d *ServiceDaemon) Start() error {
	cmdStr := fmt.Sprintf("cd %s && scripts/nebula.service start %s", d.s.dir, d.s.name)
	log.WithField("cmd", cmdStr).Debug("Try to start service")
	cmd := exec.Command("bash", "-c", cmdStr)
	err := cmd.Run()
	if err != nil {
		log.WithError(err).Errorf("Start %s failed", d.s.name)
		return err
	}
	return nil
}

func (d *ServiceDaemon) Stop() error {
	cmdStr := fmt.Sprintf("cd %s && scripts/nebula.service stop %s", d.s.dir, d.s.name)
	log.WithField("cmd", cmdStr).Debug("Try to stop service")
	cmd := exec.Command("bash", "-c", cmdStr)
	err := cmd.Run()
	if err != nil {
		log.WithError(err).Errorf("Stop %s failed", d.s.name)
		return err
	}
	return nil
}

func (d *ServiceDaemon) Status() (pb.Status, error) {
	cmdStr := fmt.Sprintf("cd %s && scripts/nebula.service status %s", d.s.dir, d.s.name)
	log.WithField("cmd", cmdStr).Debug("Try to get service's status")
	cmd := exec.Command("bash", "-c", cmdStr)
	outByte, err := cmd.Output()
	if err != nil {
		log.WithError(err).Errorf("Get status of service %s failed", d.s.name)
		return pb.Status_UNKNOWN, err
	}

	// Note: depend on the nebula scripts output now.
	outStr := string(outByte)

	// an example: [INFO] nebula-graphd(46b2aac66): Exited
	if strings.Contains(outStr, "Exit") {
		return pb.Status_EXITED, nil
	}

	// an example: [INFO] nebula-metad(46b2aac66): Running as 25859, Listening on 29559
	if strings.Contains(outStr, "Run") {
		return pb.Status_RUNNING, nil
	}

	return pb.Status_UNKNOWN, fmt.Errorf("unrecognized output: '%s'", outStr)
}
