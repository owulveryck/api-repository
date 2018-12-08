package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/owulveryck/api-repository/session"
	"github.com/owulveryck/api-repository/worker"
)

// HandleProductGet ...
func HandleProductGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	code := ps.ByName("code")
	if code == "" {
		http.Error(w, "Missing argument", http.StatusBadRequest)
		return
	}
}

type reply struct {
	ID uuid.UUID `json:"id"`
}

// ProductCreate validate the payload and stores in in the bucket
func ProductCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//ctx := r.Context()
	var products []*Product
	rdr := &io.LimitedReader{
		R: r.Body,
		N: config.MaxLength,
	}
	err := json.NewDecoder(rdr).Decode(&products)
	if err != nil {
		if rdr.N <= 0 {
			http.Error(w, fmt.Sprintf("Cannot decode elements: the payload is more than the max %v bytes", config.MaxLength), http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, fmt.Sprintf("Cannot decode elements: %v", err), http.StatusUnprocessableEntity)
		return
	}
	transactionID := uuid.New()
	t := &session.Transaction{
		ID:         transactionID.String(),
		Elements:   make([]session.Element, len(products)),
		LastUpdate: time.Now(),
	}
	for i := 0; i < len(products); i++ {
		id := products[i].ID()
		t.Elements[i] = session.Element{
			ID:     id,
			Status: http.StatusAccepted,
		}
	}
	err = session.Create(r.Context(), transactionID, t)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot create transaction: %v", err), http.StatusInternalServerError)
		return
	}

	for _, element := range products {
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
