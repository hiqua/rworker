package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"testing"

	model "github.com/hiqua/rworker/pkg"

	"github.com/go-cmd/cmd"
)

func TestCanRetrieveAddedJob(t *testing.T) {
	command := cmd.NewCmd("ls")
	uid, jid := model.UserID(123), generateJID()
	d := NewJobDirectory()
	d.updateJob(uid, jid, command)
	rcommand, err := d.retrieveJob(uid, jid)
	if err != nil {
		t.Fatal("Not able to retrieve a job that was just added.")
	}
	if rcommand != command {
		t.Fatal("Retrieved command different from the added one.")
	}
}

func TestComputeJobStatus(t *testing.T) {
	dir, err := generateTestingFolder(2)
	if err != nil {
		t.Fatal(err)
	}
	uid := model.UserID(123)
	d := NewJobDirectory()
	jid := d.AddJob(uid, "/usr/bin/ls", dir)
	time.Sleep(10 * time.Millisecond)

	status, err := d.ComputeJobStatus(uid, jid)
	if err != nil {
		t.Fatal(err)
	}

	if *status.ExitCode != 0 || status.Status != model.Done {
		t.Fatal("Exit code or status wrong")
	}

}

func TestComputeJobLog(t *testing.T) {
	dir, err := generateTestingFolder(2)
	if err != nil {
		t.Fatal(err)
	}
	uid := model.UserID(123)
	d := NewJobDirectory()
	jid := d.AddJob(uid, "/usr/bin/ls", dir)

	time.Sleep(10 * time.Millisecond)

	jobLog, err := d.ComputeJobLog(uid, jid)
	if err != nil {
		t.Fatal(err)
	}

	if jobLog.Stdout != "folder0\nfolder1" {
		t.Fatalf("Wrong output: %s\n", jobLog.Stdout)
	}
	if jobLog.Stderr != "" {
		t.Fatalf("Wrong error: %s\n", jobLog.Stderr)
	}

}

func generateTestingFolder(n int) (string, error) {
	// We don't delete the temporary folder.
	dir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		return "", err
	}
	for i := 0; i < n; i++ {
		err = os.Mkdir(dir+"/folder"+fmt.Sprint(i), 0755)
		if err != nil {
			return "", err
		}
	}
	return dir, nil
}
