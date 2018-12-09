package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/owulveryck/api-repository/object"
	"github.com/owulveryck/api-repository/session"
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
	switch r.Method {
	case "POST":
		rdr := &io.LimitedReader{
			R: r.Body,
			N: e.MaxLength,
		}
		err := json.NewDecoder(rdr).Decode(&e.Elements)
		if err != nil {
			if rdr.N <= 0 {
				http.Error(w, fmt.Sprintf("Cannot decode elements: the payload may be more than %v bytes", e.MaxLength), http.StatusRequestEntityTooLarge)
				return
			}
			http.Error(w, fmt.Sprintf("Cannot decode elements: %v", err), http.StatusUnprocessableEntity)
			return
		}
		t, u := session.NewTransaction(e.Elements)
		err = session.Create(r.Context(), u, t)
		if err != nil {
			http.Error(w, "Cannot create session: "+err.Error(), http.StatusInternalServerError)
			return
		}

		for e.Elements.Next() {
			e.JobQueue <- worker.Job{
				Payload:       e.Elements.Element(),
				TransactionID: u,
				Path:          e.Path,
			}
		}
		enc := json.NewEncoder(w)
		err = enc.Encode(reply{u})
		if err != nil {
			http.Error(w, "Cannot output transaction ID"+err.Error(), http.StatusUnprocessableEntity)
			return
		}
	default:
		http.Error(w, "Only POST method is implemented by now", http.StatusMethodNotAllowed)
		return
	}
}
