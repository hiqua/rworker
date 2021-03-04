package model

import (
	"github.com/google/uuid"
)

// JobLog contains the output / error log of the job.
type JobLog struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

// JID is the job ID type.
type JID = uuid.UUID

// UserID is the user ID type.
type UserID = int

// Status of a job (enum).
type Status = string

// The possible statuses of a job.
const (
	Running Status = "running"
	Done    Status = "done"
	Stopped Status = "stopped"
)

// JobStatus reports on all relevant information about the job (but the logs).
type JobStatus struct {
	ID        JID      `json:"id"`
	Command   string   `json:"command"`
	Arguments []string `json:"arguments"`
	Status    Status   `json:"status"`
	// Can be -1 in the stopped case.
	ExitCode *int `json:"exitCode,omitempty"`
}

// JobRequest has the fields necessary to start a job. Arguments can be empty.
type JobRequest struct {
	Command   string   `json:"command"`
	Arguments []string `json:"arguments"`
}
