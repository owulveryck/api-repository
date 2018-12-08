package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HandleAttributeCreate validate the payload and stores in in the bucket
func HandleAttributeCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	/*
		ctx := r.Context()
		dec := json.NewDecoder(r.Body)
		encResponse := json.NewEncoder(w)

		// read open bracket
		_, err := dec.Token()
		if err != nil {
			http.Error(w, "Cannot decode JSON Token: "+err.Error(), http.StatusUnprocessableEntity)
			return
		}

		// while the array contains values
		w.WriteHeader(http.StatusMultiStatus)
		var wg sync.WaitGroup
		for dec.More() {
			returnCode := ReturnCode{
				Success: true,
			}
			var a Attribute
			// decode an array value (Message)
			err := dec.Decode(&a)
			if err != nil {
				log.Println("Cannot decode Attribute: ", err)
				return
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				returnCode.Code = a.Code
				obj := bkt.Object(attributePath + a.Code)
				wobj := obj.NewWriter(ctx)
				enc := json.NewEncoder(wobj)
				err = enc.Encode(a)
				if err != nil {
					returnCode.Success = false
					returnCode.Details = append(returnCode.Details, err.Error())
					log.Println("Cannot store Attribute: ", err)
				}
				// Close, just like writing a file.
				if err := wobj.Close(); err != nil {
					returnCode.Success = false
					returnCode.Details = append(returnCode.Details, err.Error())
				}
				encResponse.Encode(returnCode)
			}()
		}

		// read closing bracket
		_, err = dec.Token()
		if err != nil {
			log.Println("Cannot decode JSON Token: ", err)
			return
		}
		wg.Wait()
	*/
}
