package utils

import (
	"strconv"
	"strings"

	"github.com/vesoft-inc/nebula-go/v3/nebula"
)

// ParseAddr parse "xx.xx.xx.xx:xx" to nebula.HostAddr
func ParseAddr(host string) (*nebula.HostAddr, error) {
	ipAddr := strings.Split(host, ":")
	port, err := strconv.ParseInt(ipAddr[1], 10, 32)
	if err != nil {
		return nil, err
	}
	return &nebula.HostAddr{ipAddr[0], nebula.Port(port)}, nil
}

func StringifyAddr(addr *nebula.HostAddr) string {
	return addr.Host + ":" + strconv.Itoa(int(addr.Port))
}
