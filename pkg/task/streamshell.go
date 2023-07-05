package task

import (
	"errors"
	"os/exec"
	"sync"

	"github.com/sirupsen/logrus"
)

type StreamShell struct {
	Shell           string
	Cmd             *exec.Cmd
	Stopped bool
}
var mu sync.RWMutex

var PipeShellMap map[string]*StreamShell 

func RunStreamShell(id string, shell string, rpcSend func(s string) error) error {
	cmd := exec.Command("bash", "-c", shell)
	pipeShell := &StreamShell{
		Shell:           shell,
		Cmd:             cmd,
		Stopped: false,
	}
	mu.Lock()
	PipeShellMap[id] = pipeShell
	mu.Unlock()
	defer ClearStreamShell(id)

	reader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	buf := make([]byte, 1024)
	
	for {
		n, err := reader.Read(buf)
		if err != nil {
			mu.RLock()
			stopped := PipeShellMap[id].Stopped
			mu.RUnlock()
			if stopped {
				return errors.New("stop stream shell")
			}
			return nil
		}
		if n > 0 {
			// rpc to push message or write to file or ignore
			err := rpcSend(string(buf[:n]))
			if err != nil {
				return err
			}
		}
	}
}

func StopStreamShell(id string) error {
	mu.Lock()
	defer mu.Unlock()

	if PipeShellMap[id] == nil {
		return nil
	}
	PipeShellMap[id].Stopped = true

	cmd := PipeShellMap[id].Cmd
	if cmd.ProcessState!=nil && cmd.ProcessState.Exited() {
		logrus.Infof("stream shell %s already stopped", id)
		return nil
	}
	return cmd.Process.Kill()
}

func ClearStreamShell(id string)error{
	mu.Lock()
	defer mu.Unlock()
	if PipeShellMap[id] == nil {
		return nil
	}
	delete(PipeShellMap, id)
	return nil
}