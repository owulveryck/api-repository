package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/owulveryck/api-repository/object"
	"github.com/owulveryck/api-repository/worker"
)

// BundlePost ...
type BundlePost struct {
	Elements  object.Iterator
	Path      string
	MaxLength int64
	JobQueue  chan<- worker.Job
}

func (e BundlePost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is implemented by now", http.StatusMethodNotAllowed)
		return
	}
	rdr := &io.LimitedReader{
		R: r.Body,
		N: e.MaxLength,
	}
	err := json.NewDecoder(rdr).Decode(&e.Elements)
	if err != nil {
		if rdr.N <= 0 {
			http.Error(w, fmt.Sprintf("Cannot decode elements: the payload may be more than %v bytes", e.MaxLength), http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, fmt.Sprintf("Cannot decode elements: %v", err), http.StatusUnprocessableEntity)
		return
	}
	transactionID := uuid.New()

	for e.Elements.Next() {
		e.JobQueue <- worker.Job{
			Payload:       e.Elements.Element(),
			TransactionID: transactionID,
			Path:          e.Path,
		}
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(reply{transactionID})
	if err != nil {
		http.Error(w, "Cannot output transaction ID"+err.Error(), http.StatusUnprocessableEntity)
		return
	}
}
