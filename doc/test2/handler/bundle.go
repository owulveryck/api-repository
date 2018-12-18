package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/owulveryck/api-repository/object"
	"github.com/owulveryck/api-repository/worker"
)

// BundlePost ...
// START_HANDLER OMIT
type BundlePost struct {
	Elements  object.Iterator
	Path      string
	MaxLength int64             // OMIT
	JobQueue  chan<- worker.Job // HL
}

// END_HANDLER OMIT

// START_SERVER OMIT
func (e BundlePost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// ...
	switch r.Method { // OMIT
	case "POST": // OMIT
		rdr := &io.LimitedReader{ // OMIT
			R: r.Body,      // OMIT
			N: e.MaxLength, // OMIT
		} // OMIT
		err := json.NewDecoder(rdr).Decode(&e.Elements)
		if err != nil { // OMIT
			if rdr.N <= 0 { // OMIT
				http.Error(w, fmt.Sprintf("Cannot decode elements: the payload may be more than %v bytes", e.MaxLength), http.StatusRequestEntityTooLarge) // OMIT
				return                                                                                                                                     // OMIT
			} // OMIT
			http.Error(w, fmt.Sprintf("Cannot decode elements: %v", err), http.StatusUnprocessableEntity)
			return
		}
		for e.Elements.Next() {
			e.JobQueue <- worker.Job{ // HL
				Payload: e.Elements.Element(), // HL
				Path:    e.Path,               // HL
			} // HL
		}
		fmt.Fprintf(w, "ok")
	default: // OMIT
		http.Error(w, "Only POST method is implemented by now", http.StatusMethodNotAllowed) // OMIT
		return                                                                               // OMIT
	}
}

// END_SERVER OMIT
