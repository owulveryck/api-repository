package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/owulveryck/api-repository/dao"
	"github.com/owulveryck/api-repository/object"
	"github.com/owulveryck/api-repository/worker"
)

// BundlePost ...
// START_HANDLER OMIT
type BundlePost struct {
	Elements  object.Iterator
	Path      string
	MaxLength int64             // OMIT
	JobQueue  chan<- worker.Job // OMIT
}

func (e BundlePost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// ...
	switch r.Method { // OMIT
	case "POST": //OMIT
		err := json.NewDecoder(r.Body).Decode(&e.Elements)
		if err != nil {
			http.Error(w, fmt.Sprintf("Cannot decode elements: %v", err), http.StatusUnprocessableEntity)
			return
		}
		for e.Elements.Next() { // HL
			err = dao.Save(r.Context(), e.Elements.Element(), e.Path)
			if err != nil {
				http.Error(w, "Cannot save "+err.Error(), http.StatusUnprocessableEntity)
				return
			}
		}
		fmt.Fprintf(w, "OK")
	default: // OMIT
		http.Error(w, "Only POST method is implemented by now", http.StatusMethodNotAllowed) // OMIT
		return                                                                               // OMIT
	} // OMIT
}

// END_HANDLER OMIT
