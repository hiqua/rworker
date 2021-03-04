# Library and implementation
The server processes requests as they come, and keep the job ids and the logs in memory until
shutdown.

* no job ordering
* no resource management
* global map to retrieve information about the processes and their output (exit code, stdout, stderr)
* the job runs with the same permissions as the server

# Authentication

* mTLS
* TLS 1.3 using one of these cipher suites: https://golang.org/pkg/crypto/tls/#pkg-constants
* certificates committed to the repository
* every certificate contains a UUID in the CN field authentifying the user

# Authorization
* each job is given a user UUID when it is created
* each client only has access to the jobs it created

# API (server and client)

For more details about the API, see [API.yaml](API.yaml).

| Method | Path       | Description                                     |
|--------|------------|-------------------------------------------------|
| GET    | /log/{id}  | Retrieve stdout / stdout as two separate fields |
| GET    | /job/{id}  | Retrieve status, incl. exit code if applicable  |
| POST   | /job       | Create a new job                                |
| DELETE | /stop/{id} | Stop a job                                      |

# CLI
* not meant to be used in scripts but just for manual usage (alternative to
    curl)
* thin wrapper around the API
* specify the certificates used manually at every interaction
* the output of the CLI is the json returned by the server
* exit value != 0 if there was an `err` field in the json from the server
* commands
    * log JOB_ID
    * status JOB_ID
    * new COMMAND [ARGUMENT]...
    * stop JOB_ID

# Trust model

We trust GitHub, meaning that we assume that the artifacts downloaded from
GitHub are trustworthy without resorting to another authenticated channel.

Every client has full access. The test setup will be one client and one server.
There is no separation between clients.

# GitHub actions

* go fmt should have no effect (code is in proper format)
* run all tests on every PR
* run all tests on every commit on master
* build artifacts on every tag

# Install script
* replace the systemd unit at install
* restarts the server after install
* unit restarts server on most failures
