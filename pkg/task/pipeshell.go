package task

import (
	"fmt"
	"os/exec"
)

type PipeShell struct {
	Shell           string
	Cmd             *exec.Cmd
	PushMessageFlag bool
}

var PipeShellMap map[int32]*PipeShell
var PipeShellId int32

func RunPipeShell(id int32, shell string) error {
	cmd := exec.Command("bash", "-c", shell)
	pipeShell := &PipeShell{
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
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := reader.Read(buf)
			if err != nil || !pipeShell.PushMessageFlag {
				return
			}
			if n > 0 {
				// todo: rpc to push message or write to file or ignore
				fmt.Println(string(buf[:n]))
			}
		}
	}()
	return nil
}

func StopPipeShell(id int32) {
	PipeShellMap[id].PushMessageFlag = false
}
