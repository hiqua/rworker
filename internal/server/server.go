package server

// TODO: use proper HTTP status codes

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"net/http"

	"github.com/google/uuid"
	"github.com/hiqua/rworker/internal/directory"

	"github.com/julienschmidt/httprouter"
)

// Serve is the main function
func Serve(address, cert, key string, knownClients ...string) error {
	log.Println("Starting server...")

	router := httprouter.New()

	jobDirectory := directory.NewJobDirectory()

	router.GET("/log/:id", fetchLog(&jobDirectory))
	router.GET("/job/:id", fetchJobStatus(&jobDirectory))

	router.POST("/job", addJob(&jobDirectory))
	router.DELETE("/stop/:id", stopJob(&jobDirectory))

	// TODO: better location for these files?
	return listenAndServeTLS(address, cert, key, router, knownClients...)
}

func listenAndServeTLS(addr, certFile, keyFile string, handler http.Handler, knownClients ...string) error {
	server := &http.Server{Addr: addr, Handler: handler}

	// TODO: do not hardcode location
	certPool, err := createClientCertPool(knownClients...)
	if err != nil {
		log.Println("Error while creating the client certificate pool.")
		return err
	}
	server.TLSConfig = &tls.Config{
		MinVersion: tls.VersionTLS13,
		ClientCAs:  certPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	return server.ListenAndServeTLS(certFile, keyFile)
}

func fetchJobStatus(jobDirectory *directory.JobDirectory) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		uid, err := extractUserID(r)
		if err != nil {
			sendError(w, err)
			return
		}

		logOperation("fetchJobStatus", uid)

		jid, err := parseJID(ps)
		if err != nil {
			sendError(w, err)
			return
		}

		jobStatus, err := jobDirectory.ComputeJobStatus(directory.NewJob(uid, jid))
		if err != nil {
			sendError(w, err)
			return
		}

		sendJSON(w, *jobStatus)
	}
}

func fetchLog(jobDirectory *directory.JobDirectory) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		uid, err := extractUserID(r)
		if err != nil {
			sendError(w, err)
			return
		}
		logOperation("fetchLog", uid)

		jid, err := parseJID(ps)
		if err != nil {
			sendError(w, err)
			return
		}

		jobLog, err := jobDirectory.ComputeJobLog(directory.NewJob(uid, jid))
		if err != nil {
			sendError(w, err)
			return
		}

		sendJSON(w, *jobLog)
	}
}

func addJob(jobDirectory *directory.JobDirectory) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		uid, err := extractUserID(r)
		if err != nil {
			sendError(w, err)
			return
		}
		logOperation("addJob", uid)

		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error while reading body: %s\n", err)
			sendError(w, err)
			return
		}

		jobRequest := directory.JobRequest{}
		err = json.Unmarshal(bs, &jobRequest)
		if err != nil {
			log.Printf("Error while parsing jobRequest: %s\n.", err)
			sendError(w, err)
			return
		}

		jid := jobDirectory.AddJob(uid, jobRequest.Command, jobRequest.Arguments...)
		sendSuccess(w, jid)
	}
}

func stopJob(jobDirectory *directory.JobDirectory) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		uid, err := extractUserID(r)
		if err != nil {
			sendError(w, err)
			return
		}
		logOperation("stopJob", uid)

		jid, err := parseJID(ps)
		if err != nil {
			sendError(w, err)
			return
		}

		err = jobDirectory.StopJob(directory.NewJob(uid, jid))
		if err != nil {
			sendError(w, err)
			return
		}
	}
}

func sendError(w http.ResponseWriter, err error) {
	var httpStatus int
	// TODO: could distinguish further
	switch err.(type) {
	case *directory.JobNotFoundError:
		httpStatus = http.StatusNotFound
	default:
		httpStatus = http.StatusInternalServerError
	}
	log.Printf("Returning an error: %s\n", err)
	sendJSONWithStatus(w, map[string]interface{}{"err": err.Error()}, httpStatus)
}

func sendSuccess(w http.ResponseWriter, jid uuid.UUID) {
	sendJSON(w, map[string]interface{}{"id": jid})
}

func sendJSON(w http.ResponseWriter, jsonObject interface{}) {
	sendJSONWithStatus(w, jsonObject, http.StatusOK)
}

func sendJSONWithStatus(w http.ResponseWriter, jsonObject interface{}, httpStatus int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bs, err := json.MarshalIndent(jsonObject, "", "  ")
	if err != nil {
		bs, err = json.Marshal(map[string]string{"err": err.Error()})
		if err != nil {
			log.Println(err)
			http.Error(w, "Unexpected error: could not marshal err.Error()", http.StatusInternalServerError)
			return
		}
		http.Error(w, string(bs), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(httpStatus)
	fmt.Fprintf(w, "%s", bs)
}

func extractUserID(r *http.Request) (uuid.UUID, error) {
	// TODO: Add custom errors for specific cases (missing CN...)
	id, err := uuid.Parse(r.TLS.PeerCertificates[0].Subject.CommonName)
	if err != nil {
		log.Println("Error while extracting the user id.")
	}
	return id, err
}

func parseJID(ps httprouter.Params) (uuid.UUID, error) {
	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		log.Println("Error while parsing the job id.")
	}
	return id, err
}

func logOperation(operation string, uid uuid.UUID) {
	log.Printf("uid: %s, op: %s", fmt.Sprint(uid), operation)
}

func createClientCertPool(certPaths ...string) (*x509.CertPool, error) {
	caCertPool := x509.NewCertPool()
	var finalErr error = nil

	for _, path := range certPaths {
		cert, err := ioutil.ReadFile(path)
		if err != nil {
			if finalErr == nil {
				finalErr = errors.New("more than one error while creating client certificate pool")
			}
			finalErr = fmt.Errorf("%w; %s", finalErr, err)
		} else {
			caCertPool.AppendCertsFromPEM(cert)
		}
	}

	return caCertPool, finalErr

}
