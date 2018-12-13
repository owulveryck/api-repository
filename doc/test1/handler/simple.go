// Package handler for documentation purpose
// START_PACKAGE OMIT
package handler

// END_PACKAGE OMIT

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/owulveryck/api-repository/dao"
	"github.com/owulveryck/api-repository/object"
)

// SimplePost handler for an IDer
// START_HANDLER OMIT
type SimplePost struct {
	Element object.IDer
	Path    string
}

func (e SimplePost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// ...
	switch r.Method { // OMIT
	case "POST": // OMIT
		err := json.NewDecoder(r.Body).Decode(&e.Element)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("Cannot decode elements: %v", err),
				http.StatusUnprocessableEntity) // HL
			return
		}
		err = dao.Save(r.Context(), e.Element, e.Path) // HL
		if err != nil {
			http.Error(w, "Cannot save "+err.Error(), http.StatusUnprocessableEntity)
			return
		}
		fmt.Fprintf(w, "OK")
	default: // OMIT
		http.Error(w, "Only POST method is implemented by now", http.StatusMethodNotAllowed) // OMIT
		return                                                                               //OMIT
	} // OMIT
	// ...
} // OMIT

// END_HANDLER OMIT
