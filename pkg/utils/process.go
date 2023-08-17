package utils

import (
	"strings"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/sirupsen/logrus"
)

func GetPidByName(name string) []int {
	name = strings.ToLower(name)
	pids := []int{}
	psList, err := process.Processes()

	if err != nil {
		return pids
	}
	for _, ps := range psList {

		process, err := ps.Cmdline()

		if err != nil {
			continue
		}
		if strings.Contains(strings.ToLower(process), name) {
			pids = append(pids, int(ps.Pid))
		}
	}

	return pids
}

func KillProcessByPids(pids []int) {
	for _, pid := range pids {
		proc, err := process.NewProcess(int32(pid))
		if err != nil {
			logrus.Error("get process failed: ", err)
			continue
		}
		if err := proc.Kill(); err != nil {
			logrus.Error("kill process failed: ", err)
		}
	}
}
