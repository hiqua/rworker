package server

import (
	"fmt"
	"strings"
	"sync"

	"github.com/go-cmd/cmd"
	"github.com/hiqua/rworker/pkg"
)

// JobDirectory provides a safe interface to the map containing all jobs for all users.
type JobDirectory struct {
	jobs map[model.UserID]map[model.JID]*cmd.Cmd
	// One global lock. Other designs could perform.
	lock sync.RWMutex
}

// NewJobDirectory initializes a new Job Directory.
func NewJobDirectory() JobDirectory {
	return JobDirectory{jobs: make(map[model.UserID]map[model.JID]*cmd.Cmd), lock: sync.RWMutex{}}

}

// AddJob add a started cmd to the given user at the returned jid
func (d *JobDirectory) AddJob(uid model.UserID, commandName string, arguments ...string) model.JID {
	command := cmd.NewCmd(commandName, arguments...)

	// TODO: we could block until the process is actually started. We could
	// then guarantee that a job said to be "running" to the client is either
	// still running or done when the client receives the answer (edge case
	// though).
	// We leave like this for now.
	command.Start()

	jid := generateJID()
	d.updateJob(uid, jid, command)
	return jid
}

// ComputeJobStatus returns the Status of the job
func (d *JobDirectory) ComputeJobStatus(uid model.UserID, jid model.JID) (*model.JobStatus, error) {
	command, err := d.retrieveJob(uid, jid)
	if err != nil {
		return nil, err
	}

	jobStatus := model.JobStatus{ID: jid, Command: command.Name, Arguments: command.Args}

	// no lock needed, .Status() safe to call concurrently
	status := command.Status()

	// See doc of the Status struct: we consider the job done if status.Error is nil, if it has run and
	// then stopped.
	if err := status.Error; err != nil {
		return nil, err
	}

	if status.StartTs > 0 && status.StopTs > 0 {
		if status.Exit == -1 {
			jobStatus.Status = model.Stopped
		} else {
			jobStatus.Status = model.Done
		}

		jobStatus.ExitCode = &status.Exit
	} else {
		jobStatus.Status = model.Running
		// jobStatus.ExitCode = nil for running jobs
	}

	return &jobStatus, nil
}

// ComputeJobLog computes the log of the job.
func (d *JobDirectory) ComputeJobLog(uid model.UserID, jid model.JID) (*model.JobLog, error) {
	command, err := d.retrieveJob(uid, jid)
	if err != nil {
		return nil, err
	}

	status := command.Status()
	stdout, stderr := strings.Join(status.Stdout, "\n"), strings.Join(status.Stderr, "\n")

	return &model.JobLog{Stdout: stdout, Stderr: stderr}, nil
}

// StopJob stops the cmd
func (d *JobDirectory) StopJob(uid model.UserID, jid model.JID) error {
	cmd, err := d.retrieveJob(uid, jid)
	if err != nil {
		return err
	}
	// .Stop() safe to call concurrently and multiple times
	// Note that only a started command can stopped, so there is a small time
	// after adding a job where, in theory, stopping it would fail.
	return cmd.Stop()
}

type UserNotFoundError struct {
	uid model.UserID
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("User id not found: %s", fmt.Sprint(e.uid))
}

type JobNotFoundError struct {
	jid model.JID
}

func (e *JobNotFoundError) Error() string {
	return fmt.Sprintf("Job id not found: %s", fmt.Sprint(e.jid))
}

// retrieveJob read the job directory safely and returns the corresponding
// command. Concurrent operations on the command are handled by the cmd library,
// which are designed to be safe.
func (d *JobDirectory) retrieveJob(uid model.UserID, jid model.JID) (*cmd.Cmd, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()

	userJobs, ok := d.jobs[uid]
	if !ok {
		return nil, &UserNotFoundError{uid}

	}
	command, ok := userJobs[jid]
	if !ok {
		return nil, &JobNotFoundError{jid}

	}
	return command, nil
}

// updateJob sets the given key (uid, jid) to the given command, safely.
func (d *JobDirectory) updateJob(uid model.UserID, jid model.JID, command *cmd.Cmd) {
	d.lock.Lock()
	defer d.lock.Unlock()

	userJobs, ok := d.jobs[uid]
	if !ok {
		userJobs = make(map[model.JID]*cmd.Cmd)
		d.jobs[uid] = userJobs
	}

	userJobs[jid] = command
}
