package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/owulveryck/api-repository/object"
	"github.com/owulveryck/api-repository/worker"
)

// SimplePost ...
type SimplePost struct {
	Element  object.IDer
	Path     string
	JobQueue chan<- worker.Job
}

func (e SimplePost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is implemented by now", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&e.Element)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot decode elements: %v", err), http.StatusUnprocessableEntity)
		return
	}
	transactionID := uuid.New()

	e.JobQueue <- worker.Job{
		Payload:       e.Element,
		TransactionID: transactionID,
		Path:          e.Path,
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(reply{transactionID})
	if err != nil {
		http.Error(w, "Cannot output transaction ID"+err.Error(), http.StatusUnprocessableEntity)
		return
	}
}
