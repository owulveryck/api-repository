package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/owulveryck/api-repository/object"
	"github.com/owulveryck/api-repository/worker"
)

func postHandler(w http.ResponseWriter, r *http.Request, elements []object.IDer, path string) {
	//ctx := r.Context()
	rdr := &io.LimitedReader{
		R: r.Body,
		N: config.MaxLength,
	}
	err := json.NewDecoder(rdr).Decode(&elements)
	if err != nil {
		if rdr.N <= 0 {
			http.Error(w, fmt.Sprintf("Cannot decode elements: the payload may be more than %v bytes", config.MaxLength), http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, fmt.Sprintf("Cannot decode elements: %v", err), http.StatusUnprocessableEntity)
		return
	}
	transactionID := uuid.New()

	for _, element := range elements {
		jobQueue <- worker.Job{
			Payload:       element,
			TransactionID: transactionID,
			Path:          productPath,
		}
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(reply{transactionID})
	if err != nil {
		http.Error(w, "Cannot output transaction ID"+err.Error(), http.StatusUnprocessableEntity)
		return
	}
}
