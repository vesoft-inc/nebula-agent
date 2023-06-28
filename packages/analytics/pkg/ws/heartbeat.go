package ws

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/sirupsen/logrus"

	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/config"
	agentConfig "github.com/vesoft-inc/nebula-agent/v3/pkg/config"
)

type MachineInfo struct {
	Host   string
	CPU    CPUInfo
	Memory MemoryInfo
	Disk   DiskInfo
}

type CPUInfo struct {
	Percentages []float64
	Avg         float64
}
type MemoryInfo struct {
	Total     uint64
	Available uint64
}

type MountPoint struct {
	Path  string
	Total uint64
	Used  uint64
}
type DiskInfo struct {
	Total       uint64
	Used        uint64
	Mountpoints []MountPoint
}

type ProcessInfo struct {
	Pid           int32
	Name          string
	CpuPercent    float64
	MemoryPercent float32
}

var timeTricker *time.Ticker

func SendHeartBeat() {
	wg := &sync.WaitGroup{}
	for _, conn := range WsClients {
		wg.Add(1)
		go func(c *websocket.Conn) {
			bytes, err := GetMachineInfo()
			if err != nil {
				logrus.Errorf("get machine info failed: %v", err)
			}
			err = c.WriteMessage(websocket.TextMessage, bytes)
			if err != nil {
				logrus.Errorf("send heartbeat to %v failed: %v", c.LocalAddr(), err)
			}
			wg.Done()
		}(conn)
	}
	wg.Wait()
	logrus.Infof("send heartbeat successfully %v", WsClients)
}

func StartHeartBeat() {
	SendHeartBeat()
	timeTricker := time.NewTicker(time.Duration(config.C.HeartBeatInterval) * time.Second)
	for range timeTricker.C {
		SendHeartBeat()
	}
}

func StopHeartBeat() {
	timeTricker.Stop()
}

func GetMachineInfo() ([]byte, error) {
	data := MachineInfo{
		CPU:    GetCPUInfo(),
		Memory: GetMemoryInfo(),
		Disk:   GetDiskInfo(),
		Host:   agentConfig.C.Agent,
	}
	text, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("get machine info failed: %v", err)
		return nil, err
	}
	return text, nil
}

func GetCPUInfo() CPUInfo {
	var cpuInfo CPUInfo
	percentages, err := cpu.Percent(time.Duration(config.C.HeartBeatInterval), false)
	if err != nil {
		logrus.Errorf("get cpu info failed: %v", err)
		return cpuInfo
	}
	var avg float64 = 0
	for _, percentage := range percentages {
		avg += percentage
	}
	avg /= float64(len(percentages))
	return CPUInfo{
		Percentages: percentages,
		Avg:         avg,
	}
}

func GetMemoryInfo() MemoryInfo {
	m, err := mem.VirtualMemory()
	if err != nil {
		logrus.Errorf("get memory info failed: %v", err)
		return MemoryInfo{}
	}

	return MemoryInfo{
		Total:     m.Total,
		Available: m.Available,
	}
}

func GetDiskInfo() DiskInfo {
	partitions, err := disk.Partitions(true)
	if err != nil {
		logrus.Errorf("get disk info failed: %v", err)
		return DiskInfo{}
	}
	var diskInfo DiskInfo = DiskInfo{
		Total:       0,
		Used:        0,
		Mountpoints: []MountPoint{},
	}
	for _, partition := range partitions {
		diskUsage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			logrus.Errorf("get disk info failed: %v", err)
			continue
		}
		diskInfo.Total += diskUsage.Total
		diskInfo.Used += diskUsage.Used
		diskInfo.Mountpoints = append(diskInfo.Mountpoints, MountPoint{
			Path:  partition.Mountpoint,
			Total: diskUsage.Total,
			Used:  diskUsage.Used,
		})
	}
	return diskInfo
}

func GetProcessInfo() []ProcessInfo {
	processes, err := process.Processes()
	if err != nil {
		logrus.Errorf("get process info failed: %v", err)
		return []ProcessInfo{}
	}
	processesInfo := []ProcessInfo{}
	for _, process := range processes {
		name, err := process.Name()
		if err != nil {
			logrus.Errorf("get process info failed: %v", err)
		}
		cpuPercent, err := process.CPUPercent()
		if err != nil {
			logrus.Errorf("get process info failed: %v", err)
			continue
		}
		memoryPercent, err := process.MemoryPercent()
		if err != nil {
			logrus.Errorf("get process info failed: %v", err)
			continue
		}
		processesInfo = append(processesInfo, ProcessInfo{
			Pid:           process.Pid,
			Name:          name,
			CpuPercent:    cpuPercent,
			MemoryPercent: memoryPercent,
		})
	}
	return processesInfo
}
