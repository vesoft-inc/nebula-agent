package types

type TaskStatus string

const (
	TaskStatusRunning     TaskStatus = "running"
	TaskStatusSuccess     TaskStatus = "success"
	TaskStatusFailed      TaskStatus = "failed"
	TaskStatusInterrupted TaskStatus = "interrupted"
	TaskStatusStopped    TaskStatus = "stopped"

)