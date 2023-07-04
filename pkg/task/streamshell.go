package task

import (
	"os/exec"
)

type StreamShell struct {
	Shell           string
	Cmd             *exec.Cmd
	PushMessageFlag bool
}

var PipeShellMap map[string]*StreamShell

func RunStreamShell(id string, shell string, rpcSend func(s string) error) error {
	cmd := exec.Command("bash", "-c", shell)
	pipeShell := &StreamShell{
		Shell:           shell,
		Cmd:             cmd,
		PushMessageFlag: true,
	}
	PipeShellMap[id] = pipeShell
	err := cmd.Start()
	if err != nil {
		return err
	}
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil || !pipeShell.PushMessageFlag {
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

func StopStreamShell(id string)error {
	PipeShellMap[id].PushMessageFlag = false
	cmd := PipeShellMap[id].Cmd
	return cmd.Process.Kill()
}
