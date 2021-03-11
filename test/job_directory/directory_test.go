package directory

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"testing"

	"github.com/google/uuid"
	"github.com/hiqua/rworker/internal/directory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComputeJobStatus(t *testing.T) {
	dir, err := generateTestingFolder(2)
	require.Nil(t, err)

	uid := uuid.New()
	d := directory.NewJobDirectory()
	jid := d.AddJob(uid, "/usr/bin/ls", dir)

	// TODO: don't sleep in tests
	time.Sleep(10 * time.Millisecond)

	status, err := d.ComputeJobStatus(directory.NewJob(uid, jid))
	require.Nil(t, err)

	assert.Equal(t, status.Status, directory.Done)
	if assert.NotNil(t, status.ExitCode) {
		assert.Equal(t, *status.ExitCode, 0)
	}
}

func TestComputeJobLog(t *testing.T) {
	dir, err := generateTestingFolder(2)
	require.Nil(t, err)

	uid := uuid.New()
	d := directory.NewJobDirectory()
	jid := d.AddJob(uid, "/usr/bin/ls", dir)

	// TODO: don't sleep in tests
	time.Sleep(10 * time.Millisecond)

	jobLog, err := d.ComputeJobLog(directory.NewJob(uid, jid))
	require.Nil(t, err)
	assert.Equal(t, jobLog.Stdout, "folder0\nfolder1")
	assert.Equal(t, jobLog.Stderr, "")
}

// TODO: could use testify/suite
func generateTestingFolder(n int) (string, error) {
	// We don't delete the temporary folder.
	dir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		return "", err
	}
	for i := 0; i < n; i++ {
		path := filepath.Join(dir, "folder"+fmt.Sprint(i))
		err = os.Mkdir(path, 0755)
		if err != nil {
			return "", err
		}
	}
	return dir, nil
}
