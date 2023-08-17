package utils

import (
	"os/exec"
	"testing"
)

func TestGetPidByName(t *testing.T) {
	// Start a process to test with
	cmd := exec.Command("sleep", "10")
	err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start process: %v", err)
	}
	defer cmd.Process.Kill()

	// Get the PID of the process
	pid := cmd.Process.Pid

	// Test the function
	pids := GetPidByName("sleep")
	if len(pids) != 1 || pids[0] != pid {
		t.Errorf("GetPidByName returned incorrect result: %v", pids)
	}
}

func TestKillProcessByPids(t *testing.T) {
	// Start a process to test with
	cmd := exec.Command("sleep", "10")
	err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start process: %v", err)
	}

	// Get the PID of the process
	pid := cmd.Process.Pid

	// Test the function
	KillProcessByPids([]int{pid})

	// Wait for the process to exit
	err = cmd.Wait()
	if err == nil {
		t.Errorf("KillProcessByPids failed to kill process")
	}
}
