package ws

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/sirupsen/logrus"

	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/clients"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/config"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/license"
	"github.com/vesoft-inc/nebula-agent/v3/packages/analytics/pkg/types"

	agentConfig "github.com/vesoft-inc/nebula-agent/v3/pkg/config"
)

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
	Free        uint64
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
	for wsHost, conn := range clients.WsClients {
		wg.Add(1)
		go func(c *websocket.Conn, wsHost string) {
			info := GetMachineInfo()
			err := c.WriteJSON(info)
			if err != nil {
				logrus.Errorf("send heartbeat to %v failed: %v", wsHost, err)
				reconnect(wsHost)
			} else {
				logrus.Infof("send heartbeat successfully to %s, status: %s:%s", c.RemoteAddr().String(), info.Body.Content["status"], info.Body.Content["statusReason"])
			}
			wg.Done()
		}(conn, wsHost)
	}
	wg.Wait()
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

func GetMachineInfo() types.Ws_Message {
	status, err := GetStatus()
	res := types.Ws_Message{
		Header: types.Ws_Message_Header{
			SendTime: time.Now().UnixMilli(),
		},
		Body: types.Ws_Message_Body{
			MsgType: types.Ws_Message_Type_Machine_Info,
			Content: map[string]interface{}{
				"status":       status,
				"statusReason": err,
				"cpu":          GetCPUInfo(),
				"memory":       GetMemoryInfo(),
				"agent":        agentConfig.C.Agent,
			},
		},
	}
	return res
}

func GetStatus() (string, error) {
	algoPath := path.Join(config.C.AnalyticsPath, "scripts/run_algo.sh")
	if _, err := os.Stat(algoPath); os.IsNotExist(err) {
		return "error", fmt.Errorf("load lm config failed: %v", err)
	}
	if license.LicenseApplyHandler.IsStopService {
		return "error", fmt.Errorf(license.LicenseApplyHandler.Error)
	}
	return "ok", nil
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

func SendAgentChangeToExplorer() {
	explorerHosts := make([]string, 0, len(clients.WsClients))
	for k := range clients.WsClients {
		explorerHosts = append(explorerHosts, k)
	}
	for _, conn := range clients.WsClients {
		res := types.Ws_Message{
			Header: types.Ws_Message_Header{
				SendTime: time.Now().UnixMilli(),
			},
			Body: types.Ws_Message_Body{
				MsgType: types.Ws_Message_Type_Agent,
				Content: map[string]interface{}{
					"activeExplorerHosts": explorerHosts,
					"agent":               agentConfig.C.Agent,
					"explorerHosts":       config.C.ExplorerHosts,
				},
			},
		}
		bytes, err := json.Marshal(res)
		if err != nil {
			logrus.Errorf("marshal message error: %v", err)
			continue
		}
		err = conn.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			logrus.Errorf("write message error: %v", err)
			continue
		}
	}
}
