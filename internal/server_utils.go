package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	model "github.com/hiqua/rworker/pkg"
	"github.com/julienschmidt/httprouter"
)

func sendError(w http.ResponseWriter, err error) {
	var httpStatus int
	// TODO: could distinguish further
	switch err.(type) {
	case *UserNotFoundError:
		httpStatus = http.StatusBadRequest
	case *JobNotFoundError:
		httpStatus = http.StatusBadRequest
	default:
		httpStatus = http.StatusInternalServerError
	}
	log.Printf("Returning an error: %s\n", err)
	sendJSONWithStatus(w, map[string]interface{}{"err": err.Error()}, httpStatus)
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

func dummyExtractUserID() model.UserID {
	// TODO: implement actual function extracting user id from client certificate
	return model.UserID(3)
}

// For debugging purposes
// func dummyGenerateJID() model.JID {
// 	jid, _ := uuid.FromBytes([]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
// 	return model.JID(jid)
// }

func parseJID(ps httprouter.Params) (model.JID, error) {
	uuid, err := uuid.Parse(ps.ByName("id"))
	return model.JID(uuid), err
}

func generateJID() model.JID {
	return model.JID(uuid.New())
}

func logOperation(operation string, uid model.UserID) {
	log.Printf("uid: %s, op: %s", fmt.Sprint(uid), operation)
}
