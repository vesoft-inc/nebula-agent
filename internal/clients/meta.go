package clients

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vesoft-inc/fbthrift/thrift/lib/go/thrift"
	"github.com/vesoft-inc/nebula-agent/v3/internal/utils"
	"github.com/vesoft-inc/nebula-go/v3/nebula"
	"github.com/vesoft-inc/nebula-go/v3/nebula/meta"
)

const (
	defaultTimeout time.Duration = 120 * time.Second
)

type MetaConfig struct {
	GitInfoSHA   string
	HBInterval   int
	MetaAddr     *nebula.HostAddr // meta service address to connect
	AgentAddr    *nebula.HostAddr // info to be reported to the meta service
	TLSConfig    *tls.Config
	HandShakeKey string
}

func NewMetaConfig(agentAddr, metaAddr, gitSHA, handshakeKey string, hbInterval int, tlsConfig *tls.Config) (*MetaConfig, error) {
	cfg := &MetaConfig{
		GitInfoSHA:   gitSHA,
		HBInterval:   hbInterval,
		TLSConfig:    tlsConfig,
		HandShakeKey: handshakeKey,
	}
	var err error

	if cfg.AgentAddr, err = utils.ParseAddr(agentAddr); err != nil {
		log.WithError(err).WithField("address", agentAddr).Info("Parse agent address failed.")
		return nil, err
	}
	if cfg.MetaAddr, err = utils.ParseAddr(metaAddr); err != nil {
		log.WithError(err).WithField("address", metaAddr).Info("Parse meta service address failed.")
		return nil, err
	}

	return cfg, nil
}

// NebulaMeta is the client to communicate with nebula meta servie, such as heartbeat\getMetaInfo.
// Through which, the agent could get services to monitor dynamically
type NebulaMeta struct {
	client *meta.MetaServiceClient
	config *MetaConfig

	mu       sync.RWMutex
	services map[string]*meta.ServiceInfo
}

func NewMeta(config *MetaConfig) (*NebulaMeta, error) {
	m := &NebulaMeta{
		config: config,
	}
	var err error

	if m.client, err = connect(m.config.MetaAddr, m.config.AgentAddr, m.config.TLSConfig, m.config.HandShakeKey); err != nil {
		return nil, err
	}

	go m.heartbeatLoop()
	return m, nil
}

// Connect to meta service by given meta address
// We will reconnect when:
//  1. the meta service leader changed
//  2. get individual info from individual meta service, such as dir info
func connect(metaAddr, agentAddr *nebula.HostAddr, tlsConfig *tls.Config, handshakeKey string) (*meta.MetaServiceClient, error) {
	addr := utils.StringifyAddr(metaAddr)
	log.WithField("meta address", addr).Info("try to connect meta service")

	timeoutOption := thrift.SocketTimeout(defaultTimeout)
	addressOption := thrift.SocketAddr(addr)

	var (
		err  error
		sock thrift.Transport
	)
	if tlsConfig != nil {
		tlsConfig.MaxVersion = tls.VersionTLS12
		sock, err = thrift.NewSSLSocketTimeout(addr, tlsConfig, defaultTimeout)
	} else {
		sock, err = thrift.NewSocket(timeoutOption, addressOption)
	}
	if err != nil {
		log.WithError(err).WithField("addr", addr).Error("open socket failed")
		return nil, err
	}

	bufferedTranFactory := thrift.NewBufferedTransportFactory(128 << 10)
	transport := thrift.NewFramedTransport(bufferedTranFactory.GetTransport(sock))
	pf := thrift.NewBinaryProtocolFactoryDefault()
	client := meta.NewMetaServiceClientFactory(transport, pf)
	if err := client.CC.Open(); err != nil {
		log.WithError(err).WithField("addr", addr).Error("open meta failed")
		return nil, err
	}

	req := newVerifyClientVersionReq(agentAddr, handshakeKey)
	resp, err := client.VerifyClientVersion(req)
	if err != nil || resp.Code != nebula.ErrorCode_SUCCEEDED {
		if err == nil {
			err = fmt.Errorf("verifyClientVersion err: %s", string(resp.GetErrorMsg()))
		}

		log.WithError(err).WithField("addr", metaAddr).Error(string(resp.GetErrorMsg()))
		client.Close()

		return nil, err
	}

	log.WithField("meta address", addr).Info("connect meta server successfully")
	return client, nil
}

func (m *NebulaMeta) heartbeatLoop() {
	err := m.heartbeat()
	if err != nil {
		log.WithError(err).Error("Heartbeat to meta failed")
	}

	t := time.NewTicker(time.Duration(m.config.HBInterval) * time.Second)
	for {
		select {
		case <-t.C:
			err := m.heartbeat()
			if err != nil {
				log.WithError(err).Error("Heartbeat to meta failed")
			}
		}
	}
}

func (m *NebulaMeta) heartbeat() error {
	req := &meta.AgentHBReq{
		Host:       m.config.AgentAddr,
		GitInfoSha: []byte(m.config.GitInfoSHA),
	}

	try := 1
	// retry only when leader change
	for {
		resp, err := m.client.AgentHeartbeat(req)
		if err != nil {
			c, err := connect(m.config.MetaAddr, m.config.AgentAddr, m.config.TLSConfig, m.config.HandShakeKey)
			if err != nil {
				return err
			}
			m.client = c

			log.WithError(err).Error("heartbeat to nebula meta service failed")
			return err
		}

		if resp.GetCode() == nebula.ErrorCode_SUCCEEDED {
			m.refreshInfo(resp.GetServiceList())

			m.mu.RLock()
			for _, s := range m.services {
				log.WithField("service", utils.StringifyService(s)).Debug("Get service info from heartbeat.")
			}
			m.mu.RUnlock()
			return nil
		}

		if resp.GetCode() == nebula.ErrorCode_E_LEADER_CHANGED {
			if resp.GetLeader() == meta.ExecResp_Leader_DEFAULT {
				return LeaderNotFoundError
			}
			m.config.MetaAddr = resp.GetLeader()
			c, err := connect(m.config.MetaAddr, m.config.AgentAddr, m.config.TLSConfig, m.config.HandShakeKey)
			if err != nil {
				return err
			}
			m.client = c

			log.Infof("retry heatbeat, meta leader changed, try times=%d.", try)
			time.Sleep(time.Duration(try) * time.Second)
			try *= 2
			continue
		}

		log.WithField("error", resp.GetCode().String()).Error("Heartbeat to nebula failed.")
		return fmt.Errorf("heartbeat to nebula failed: %s", resp.GetCode().String())
	}
}

// getMetaDirInfo get individual meta dir info of given meta service.
// Because follower meta service could not report its dir info to the leader one.
func getMetaDirInfo(metaAddr, agentAddr *nebula.HostAddr, tlsConfig *tls.Config, handshakeKey string) (*nebula.DirInfo, error) {
	log.WithField("addr", utils.StringifyAddr(metaAddr)).
		Debugf("Try to get dir info from meta service: %s\n", utils.StringifyAddr(metaAddr))
	c, err := connect(metaAddr, agentAddr, tlsConfig, handshakeKey)
	if err != nil {
		return nil, err
	}

	defer func() {
		e := c.Close()
		if e != nil {
			log.WithError(e).WithField("host", metaAddr.String()).Error("Close error when get meta dir info.")
		}
	}()

	req := &meta.GetMetaDirInfoReq{}
	resp, err := c.GetMetaDirInfo(req)
	if err != nil {
		return nil, err
	}

	if resp.GetCode() != nebula.ErrorCode_SUCCEEDED {
		return nil, fmt.Errorf("get meta dir info failed: %v", resp.GetCode())
	}

	return resp.GetDir(), nil
}

// refreshInfo refresh the meta service's dir info by rpc
func (m *NebulaMeta) refreshInfo(services []*meta.ServiceInfo) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	newServices := make(map[string]*meta.ServiceInfo)
	for _, s := range services {
		k := utils.StringifyAddr(s.GetAddr())
		if s.GetRole() == meta.HostRole_META {
			d, err := getMetaDirInfo(s.Addr, m.config.AgentAddr, m.config.TLSConfig, m.config.HandShakeKey)
			if err != nil {
				return fmt.Errorf("get meta dir for %s failed: %w", k, err)
			}
			log.WithField("root", string(d.Root)).WithField("data", string(d.Data[0])).
				Debugf("Get dir info from meta: %s successfully.", utils.StringifyAddr(s.Addr))
			s.Dir = d
		}

		newServices[k] = s
	}
	m.services = newServices

	return nil
}

func newVerifyClientVersionReq(agentAddr *nebula.HostAddr, handshakeKey string) *meta.VerifyClientVersionReq {
	if handshakeKey == "" {
		handshakeKey = nebula.Version
	}

	return &meta.VerifyClientVersionReq{
		ClientVersion: []byte(handshakeKey),
		Host:          agentAddr,
	}
}
