package directory

import (
	"fmt"
	"strings"
	"sync"

	"github.com/go-cmd/cmd"
	"github.com/google/uuid"
)

// JobLog contains the output / error log of the job.
type JobLog struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

// Status of a job (enum).
type Status string

// The possible statuses of a job.
const (
	Running Status = "running"
	Done    Status = "done"
	Stopped Status = "stopped"
)

// JobStatus reports on all relevant information about the job (but the logs).
type JobStatus struct {
	ID        uuid.UUID `json:"id"`
	Command   string    `json:"command"`
	Arguments []string  `json:"arguments"`
	Status    Status    `json:"status"`
	// Can be -1 in the stopped case.
	ExitCode *int `json:"exitCode,omitempty"`
}

// JobRequest has the fields necessary to start a job. Arguments can be empty.
type JobRequest struct {
	Command   string   `json:"command"`
	Arguments []string `json:"arguments"`
}

// JobID represents the id of a job (response to an addJob operation).
type JobID struct {
	ID uuid.UUID
}

// JobDirectory provides a safe interface to the map containing all jobs for all users.
type JobDirectory struct {
	jobs map[Job]*cmd.Cmd
	// One global lock. Other designs could perform.
	lock sync.RWMutex
}

// Job is the key to a cmd run by a job.
type Job struct {
	userID uuid.UUID
	jobID  uuid.UUID
}

func NewJob(userID, jobID uuid.UUID) Job {
	return Job{userID, jobID}
}

// NewJobDirectory initializes a new Job Directory.
func NewJobDirectory() JobDirectory {
	return JobDirectory{jobs: make(map[Job]*cmd.Cmd), lock: sync.RWMutex{}}

}

// AddJob add a started cmd to the given user at the returned jid
func (d *JobDirectory) AddJob(uid uuid.UUID, commandName string, arguments ...string) uuid.UUID {
	command := cmd.NewCmd(commandName, arguments...)

	// TODO: we could block until the process is actually started. We could
	// then guarantee that a job said to be "running" to the client is either
	// still running or done when the client receives the answer (edge case
	// though).
	// We leave like this for now.
	command.Start()

	jid := uuid.New()
	d.updateJob(Job{uid, jid}, command)
	return jid
}

// ComputeJobStatus returns the Status of the job
func (d *JobDirectory) ComputeJobStatus(job Job) (*JobStatus, error) {
	command, err := d.retrieveJob(job)
	if err != nil {
		return nil, err
	}

	jobStatus := JobStatus{ID: job.jobID, Command: command.Name, Arguments: command.Args}

	// no lock needed, .Status() safe to call concurrently
	status := command.Status()

	// See doc of the Status struct: we consider the job done if status.Error is nil, if it has run and
	// then stopped.
	if err := status.Error; err != nil {
		return nil, err
	}

	if status.StartTs > 0 && status.StopTs > 0 {
		if status.Exit == -1 {
			jobStatus.Status = Stopped
		} else {
			jobStatus.Status = Done
		}

		jobStatus.ExitCode = &status.Exit
	} else {
		jobStatus.Status = Running
		// jobStatus.ExitCode = nil for running jobs
	}

	return &jobStatus, nil
}

// ComputeJobLog computes the log of the job.
func (d *JobDirectory) ComputeJobLog(job Job) (*JobLog, error) {
	command, err := d.retrieveJob(job)
	if err != nil {
		return nil, err
	}

	status := command.Status()
	stdout, stderr := strings.Join(status.Stdout, "\n"), strings.Join(status.Stderr, "\n")

	return &JobLog{Stdout: stdout, Stderr: stderr}, nil
}

// StopJob stops the cmd
func (d *JobDirectory) StopJob(job Job) error {
	cmd, err := d.retrieveJob(job)
	if err != nil {
		return err
	}
	// .Stop() safe to call concurrently and multiple times
	// Note that only a started command can stopped, so there is a small time
	// after adding a job where, in theory, stopping it would fail.
	return cmd.Stop()
}

// JobNotFoundError is returned when a job could not be found in the directory.
type JobNotFoundError struct {
	job Job
}

func (e *JobNotFoundError) Error() string {
	return fmt.Sprintf("Job not found: %s", fmt.Sprint(e.job.jobID))
}

// retrieveJob read the job directory safely and returns the corresponding
// command. Concurrent operations on the command are handled by the cmd library,
// which are designed to be safe.
func (d *JobDirectory) retrieveJob(job Job) (*cmd.Cmd, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()

	command, ok := d.jobs[job]
	if !ok {
		return nil, &JobNotFoundError{job}

	}
	return command, nil
}

// updateJob sets the given key (uid, jid) to the given command, safely.
func (d *JobDirectory) updateJob(job Job, command *cmd.Cmd) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.jobs[job] = command
}
