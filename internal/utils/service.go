package utils

import (
	"fmt"
	"strings"

	"github.com/vesoft-inc/nebula-go/v3/nebula/meta"
)

func StringifyService(s *meta.ServiceInfo) string {
	m := make(map[string]string)
	m["addr"] = StringifyAddr(s.GetAddr())
	m["role"] = s.GetRole().String()
	m["root_dir"] = string(s.GetDir().GetRoot())

	dataDirs := make([]string, 0, len(s.GetDir().GetData()))
	for _, d := range s.GetDir().GetData() {
		if len(d) != 0 {
			dataDirs = append(dataDirs, string(d))
		}
	}
	m["data_dirs"] = strings.Join(dataDirs, ",")

	return fmt.Sprintf("%v", m)
}
