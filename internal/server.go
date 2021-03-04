package server

// TODO: use proper HTTP status codes

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	model "github.com/hiqua/rworker/pkg"

	"github.com/julienschmidt/httprouter"
)

// Serve is the main function
func Serve() error {
	log.Println("Starting server...")

	router := httprouter.New()

	jobDirectory := NewJobDirectory()

	router.GET("/log/:id", fetchLog(&jobDirectory))
	router.GET("/job/:id", fetchJobStatus(&jobDirectory))

	router.POST("/job", addJob(&jobDirectory))
	router.DELETE("/stop", stopJob(&jobDirectory))

	// TODO: better location for these files?
	// TODO: listen to 0.0.0.0 once mTLS setup
	return listenAndServeTLS("127.0.0.1:8443", "certs/server/cert.pem", "certs/server/key.pem", router)
}

// Taken from the http package
// TODO: maybe there's a better way to enforce TLS1.3 rather than rewrite this function?
func listenAndServeTLS(addr, certFile, keyFile string, handler http.Handler) error {
	server := &http.Server{Addr: addr, Handler: handler}
	server.TLSConfig = &tls.Config{
		MinVersion: tls.VersionTLS13,
	}

	return server.ListenAndServeTLS(certFile, keyFile)
}

func fetchJobStatus(jobDirectory *JobDirectory) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		uid := dummyExtractUserID()
		logOperation("fetchJobStatus", uid)

		jid, err := parseJID(ps)
		if err != nil {
			sendError(w, err)
			return
		}

		jobStatus, err := jobDirectory.ComputeJobStatus(uid, jid)
		if err != nil {
			sendError(w, err)
			return
		}

		sendJSON(w, *jobStatus)
	}
}

func fetchLog(jobDirectory *JobDirectory) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		uid := dummyExtractUserID()
		logOperation("fetchLog", uid)

		jid, err := parseJID(ps)
		if err != nil {
			sendError(w, err)
			return
		}

		jobLog, err := jobDirectory.ComputeJobLog(uid, jid)
		if err != nil {
			sendError(w, err)
			return
		}

		sendJSON(w, *jobLog)
	}
}

func addJob(jobDirectory *JobDirectory) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		uid := dummyExtractUserID()
		logOperation("addJob", uid)

		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error while reading body.")
			sendError(w, err)
			return
		}

		jobRequest := model.JobRequest{}
		err = json.Unmarshal(bs, &jobRequest)
		if err != nil {
			log.Printf("Error while parsing jobRequest: %s\n.", err)
			sendError(w, err)
			return
		}

		jid := jobDirectory.AddJob(uid, jobRequest.Command, jobRequest.Arguments...)
		sendJSON(w, map[string]interface{}{"id": jid})
	}
}

func stopJob(jobDirectory *JobDirectory) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		uid := dummyExtractUserID()
		logOperation("stopJob", uid)

		jid, err := parseJID(ps)
		if err != nil {
			sendError(w, err)
			return
		}

		err = jobDirectory.StopJob(uid, jid)
		if err != nil {
			sendError(w, err)
			return
		}
	}
}
