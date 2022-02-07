package client

import (
	"github.com/vesoft-inc/nebula-go/v3/nebula"
)

// Config is the agent client config
type Config struct {
	Addr *nebula.HostAddr
}
